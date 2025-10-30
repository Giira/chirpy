package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileServerHits atomic.Int32
}

func main() {

	port := "8080"
	root := "."
	cfg := &apiConfig{}

	serveMux := http.NewServeMux()
	file_handler := http.StripPrefix("/app", http.FileServer(http.Dir(root)))

	serveMux.Handle("/app/", cfg.metricInc(file_handler))

	serveMux.HandleFunc("GET /api/healthz", handleReady)
	serveMux.HandleFunc("GET /admin/metrics", cfg.handleHits)

	serveMux.HandleFunc("POST /admin/reset", cfg.handleReset)
	serveMux.HandleFunc("POST /api/validate_chirp", handleValidity)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: serveMux,
	}

	fmt.Printf("Server active on port %v", port)
	log.Fatal(server.ListenAndServe())

}
