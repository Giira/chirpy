package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/Giira/chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileServerHits atomic.Int32
	queries        *database.Queries
	platform       string
}

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("error launching database: %w", err)
	}
	dbQueries := database.New(db)

	port := "8080"
	root := "."
	cfg := &apiConfig{
		queries:  dbQueries,
		platform: os.Getenv("PLATFORM"),
	}

	serveMux := http.NewServeMux()
	file_handler := http.StripPrefix("/app", http.FileServer(http.Dir(root)))

	serveMux.Handle("/app/", cfg.metricInc(file_handler))

	serveMux.HandleFunc("GET /api/healthz", handleReady)
	serveMux.HandleFunc("GET /admin/metrics", cfg.handleHits)

	serveMux.HandleFunc("POST /admin/reset", cfg.handleReset)
	serveMux.HandleFunc("POST /api/validate_chirp", handleValidity)
	serveMux.HandleFunc("POST /api/users", cfg.handleCreateUser)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: serveMux,
	}

	fmt.Printf("Server active on port %v", port)
	log.Fatal(server.ListenAndServe())

}
