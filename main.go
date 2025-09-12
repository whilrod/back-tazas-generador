package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"generadorPlantillas/db"
	"generadorPlantillas/handlers"

	"github.com/rs/cors"
)

func main() {
	// Inicializar DB
	database := db.InitDB()
	defer database.Close()

	// Registrar Handlers en un mux
	mux := http.NewServeMux()
	handlers.RegisterHandlersWithMux(mux, database) // vamos a crear esta función en handlers

	// Configurar CORS
	allowedOrigin := os.Getenv("ALLOWED_ORIGIN")
	if allowedOrigin == "" {
		allowedOrigin = "http://localhost:5173" // fallback para local
	}

	// Soporte para múltiples orígenes separados por coma
	origins := strings.Split(allowedOrigin, ",")

	c := cors.New(cors.Options{
		AllowedOrigins: origins,
		//AllowedOrigins:   []string{allowedOrigin},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: true,
	})

	// Aplicar middleware CORS
	handler := c.Handler(mux)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // fallback para local
	}
	log.Printf("🌍 Servidor escuchando en %s:%s", allowedOrigin, port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
