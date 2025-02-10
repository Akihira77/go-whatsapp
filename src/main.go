package main

import (
	"log"
	"os"

	"github.com/Akihira77/go_whatsapp/src/store"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error load .env file: %v", err)
	}

	store := store.NewStore()
	// store.Migrate()

	port := os.Getenv("PORT")
	if err := NewHttpServer(port, store); err != nil {
		log.Fatalf("Error running server: %v", err)
	}
}
