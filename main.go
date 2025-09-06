package main

import (
	"log"
	"net/http"

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
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: true,
	})

	// Aplicar middleware CORS
	handler := c.Handler(mux)

	port := ":8080"
	log.Printf("🌍 Servidor escuchando en http://localhost%s", port)
	log.Fatal(http.ListenAndServe(port, handler))
}
