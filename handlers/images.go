package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"generadorPlantillas/models"
	"generadorPlantillas/utils"

	"github.com/jung-kurt/gofpdf"
)

// ---------------------------
// Registro de Handlers
// ---------------------------
func RegisterHandlers(db *sql.DB) {
	RegisterHandlersWithMux(http.DefaultServeMux, db)
}

func RegisterHandlersWithMux(mux *http.ServeMux, db *sql.DB) {
	mux.HandleFunc("/images", func(w http.ResponseWriter, r *http.Request) {
		ListImagesHandler(db, w, r)
	})
	mux.HandleFunc("/images/hashtag", func(w http.ResponseWriter, r *http.Request) {
		SearchImagesByHashtagHandler(db, w, r)
	})
	mux.HandleFunc("/images/pdf", func(w http.ResponseWriter, r *http.Request) {
		GeneratePDFHandler(db, w, r)
	})
}

// ---------------------------
// Helpers para paginación
// ---------------------------
func getPaginationParams(r *http.Request) (page, limit, offset int) {
	// Defaults
	page = 1
	limit = 20

	if p := r.URL.Query().Get("page"); p != "" {
		if val, err := strconv.Atoi(p); err == nil && val > 0 {
			page = val
		}
	}
	if l := r.URL.Query().Get("limit"); l != "" {
		if val, err := strconv.Atoi(l); err == nil && val > 0 {
			limit = val
		}
	}

	offset = (page - 1) * limit
	return
}

// ---------------------------
// Handlers
// ---------------------------

// ListImagesHandler lista imágenes con paginación
func ListImagesHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	// valores por defecto
	page := 1
	limit := 20

	if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
		page = p
	}
	if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
		limit = l
	}

	offset := (page - 1) * limit

	// contar total de registros
	var total int
	err := db.QueryRow(`SELECT COUNT(*) FROM "imagenes"`).Scan(&total)
	if err != nil {
		log.Printf("❌ Error contando registros: %v", err)
		http.Error(w, "Error consultando la base de datos", http.StatusInternalServerError)
		return
	}

	// consulta con paginación
	query := `SELECT uuid, url_image, url_thumbnail, hashtags, xata_createdat, size_kb 
	          FROM "imagenes" 
	          ORDER BY xata_createdat DESC
	          LIMIT $1 OFFSET $2`

	rows, err := db.Query(query, limit, offset)
	if err != nil {
		log.Printf("❌ Error consultando la base de datos: %v", err)
		http.Error(w, "Error consultando la base de datos", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var results []models.Image
	for rows.Next() {
		var img models.Image
		var hashtagsRaw string

		if err := rows.Scan(&img.UUID, &img.URLImage, &img.URLThumbnail, &hashtagsRaw, &img.CreatedAt, &img.SizeKb); err != nil {
			log.Printf("❌ Error leyendo fila: %v", err)
			http.Error(w, "Error leyendo filas", http.StatusInternalServerError)
			return
		}

		img.Hashtags = utils.ParsePgArray(hashtagsRaw)
		results = append(results, img)
	}

	// total de páginas
	totalPages := (total + limit - 1) / limit

	response := map[string]interface{}{
		"results":     results,
		"page":        page,
		"limit":       limit,
		"total":       total,
		"total_pages": totalPages,
	}

	utils.RespondJSON(w, response)
}

// SearchImagesByHashtagHandler busca imágenes con hashtags y paginación
func SearchImagesByHashtagHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	tags := r.URL.Query()["tag"]

	// Si no hay tags, devolver lista normal (como si fuese /images)
	if len(tags) == 0 {
		ListImagesHandler(db, w, r)
		return
	}

	page, limit, offset := getPaginationParams(r)

	// Construir array PostgreSQL
	pgArray := "{" + strings.Join(tags, ",") + "}"

	query := `
		SELECT uuid, url_image, url_thumbnail, hashtags, xata_createdat, size_kb
		FROM imagenes
		WHERE hashtags && $1
		ORDER BY xata_createdat DESC
		LIMIT $2 OFFSET $3
	`

	images, err := queryImages(db, query, pgArray, limit, offset)
	if err != nil {
		http.Error(w, "Error consultando la base de datos", http.StatusInternalServerError)
		return
	}

	// Contar total con filtro
	var total int
	err = db.QueryRow(`SELECT COUNT(*) FROM imagenes WHERE hashtags && $1`, pgArray).Scan(&total)
	if err != nil {
		http.Error(w, "Error contando registros", http.StatusInternalServerError)
		return
	}
	totalPages := (total + limit - 1) / limit

	utils.RespondJSON(w, map[string]interface{}{
		"results":     images,
		"page":        page,
		"limit":       limit,
		"total":       total,
		"total_pages": totalPages,
	})
}

// ---------------------------
// Función interna para consultas
// ---------------------------
func queryImages(db *sql.DB, query string, args ...interface{}) ([]models.Image, error) {
	rows, err := db.Query(query, args...)
	if err != nil {
		log.Printf("❌ Error consultando la base de datos: %v", err)
		return nil, err
	}
	defer rows.Close()

	var results []models.Image
	for rows.Next() {
		var img models.Image
		var hashtagsRaw string

		if err := rows.Scan(&img.UUID, &img.URLImage, &img.URLThumbnail, &hashtagsRaw, &img.CreatedAt, &img.SizeKb); err != nil {
			log.Printf("❌ Error leyendo fila: %v", err)
			return nil, err
		}

		img.Hashtags = utils.ParsePgArray(hashtagsRaw)
		results = append(results, img)
	}

	return results, nil
}

// Request con los UUID seleccionados
type PDFRequest struct {
	UUIDs []string `json:"uuids"`
}

func GeneratePDFHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	var req PDFRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	if len(req.UUIDs) == 0 {
		http.Error(w, "No se enviaron UUIDs", http.StatusBadRequest)
		return
	}

	// Crear PDF
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetFont("Arial", "", 14)

	// Posiciones fijas Y para cada slot de imagen
	positionsY := []float64{0, 101.45, 202.9} // ajustado para que queden bien distribuidas
	imgWidth := 210.0
	imgHeight := 94.1

	for i, uuid := range req.UUIDs {
		var urlImage string
		err := db.QueryRow(`SELECT url_image FROM imagenes WHERE uuid = $1`, uuid).Scan(&urlImage)
		if err != nil {
			fmt.Printf("❌ Error obteniendo imagen %s: %v\n", uuid, err)
			continue
		}

		// Descargar imagen remota
		resp, err := http.Get(urlImage)
		if err != nil {
			fmt.Printf("❌ Error descargando imagen %s: %v\n", urlImage, err)
			continue
		}
		defer resp.Body.Close()

		// Registrar imagen en memoria
		opt := gofpdf.ImageOptions{ImageType: "PNG", ReadDpi: true}
		imgName := uuid
		pdf.RegisterImageOptionsReader(imgName, opt, resp.Body)

		// Si es la primera de la página, agregamos nueva página
		if i%3 == 0 {
			pdf.AddPage()
		}

		// Posición Y en base al índice dentro de la página
		slot := i % 3
		y := positionsY[slot]

		// Dibujar imagen ocupando el ancho completo de la página
		pdf.ImageOptions(imgName, 0, y, imgWidth, imgHeight, false, opt, 0, "")
	}

	// Enviar PDF como respuesta
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename=imagenes.pdf")
	if err := pdf.Output(w); err != nil {
		http.Error(w, "Error generando PDF", http.StatusInternalServerError)
	}
}
