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

func (cfg *apiConfig) metricInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileServerHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func main() {

	port := "8080"
	root := "."
	cfg := &apiConfig{}

	serveMux := http.NewServeMux()
	file_handler := http.StripPrefix("/app", http.FileServer(http.Dir(root)))

	serveMux.Handle("/app", cfg.metricInc(file_handler))
	serveMux.HandleFunc("/healthz", handleReady)
	serveMux.HandleFunc("/metrics", cfg.handleHits)
	serveMux.HandleFunc("/reset", cfg.handleReset)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: serveMux,
	}

	fmt.Printf("Server active on port %v", port)
	log.Fatal(server.ListenAndServe())

}
