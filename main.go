package main

import (
	"log"

	"github.com/joho/godotenv"
)

func main() {
    if err := godotenv.Load(); err != nil {
        log.Fatalf("Oops: %w", err)
    }

    if err := startServer(); err != nil {
        log.Fatalf("Oops: %w", err)
    }
}
