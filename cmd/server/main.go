package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Georgy-collab/Black-Box/internal/server"
)

func main() {
	addr := ":8000"
	if port := os.Getenv("PORT"); port != "" {
		addr = ":" + port
	}

	log.Printf("User Service listening on http://127.0.0.1%s", addr)
	log.Printf("Swagger UI: http://127.0.0.1%s/docs", addr)
	log.Printf("OpenAPI:   http://127.0.0.1%s/openapi.json", addr)

	if err := http.ListenAndServe(addr, server.New()); err != nil {
		log.Fatal(err)
	}
}
