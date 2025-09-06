package db

import (
	"database/sql"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func InitDB() *sql.DB {
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️ No se encontró archivo .env, usando variables del sistema")
	}

	dbURL := os.Getenv("XATA_DATABASE_URL")
	if dbURL == "" {
		log.Fatal("❌ No se encontró la variable de entorno XATA_DATABASE_URL")
	}

	database, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("❌ Error abriendo conexión: %v", err)
	}

	if err := database.Ping(); err != nil {
		log.Fatalf("❌ Error conectando a la base de datos: %v", err)
	}

	log.Println("✅ Conexión exitosa con Xata 🚀")
	return database
}
