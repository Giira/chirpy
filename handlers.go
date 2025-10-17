package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

func handleReady(writer http.ResponseWriter, req *http.Request) {
	req.Header.Set("Content-Type", "text/plain; charset=utf-8")
	writer.WriteHeader(200)
	writer.Write([]byte("OK"))
}

func (cfg *apiConfig) handleHits(writer http.ResponseWriter, req *http.Request) {
	req.Header.Set("Content-Type", "text/html; charset=utf-8")
	writer.WriteHeader(200)
	text := fmt.Sprintf("<html><body><h1>Welcome, Chirpy Admin</h1><p>Chirpy has been visited %d times!</p></body></html>", cfg.fileServerHits.Load())
	writer.Write([]byte(text))
}

func (cfg *apiConfig) handleReset(writer http.ResponseWriter, req *http.Request) {
	req.Header.Set("Content-Type", "text/plain; charset=utf-8")
	writer.WriteHeader(200)
	before := cfg.fileServerHits.Load()
	cfg.fileServerHits.Store(0)
	after := cfg.fileServerHits.Load()
	text := fmt.Sprintf("Hits: %v\nHits reset\nHits: %v", before, after)
	writer.Write([]byte(text))

}

func handleValidity(writer http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		writer.WriteHeader(500)
		return
	}

	type returnVals struct {
		CreatedAt time.Time `json:"created_at"`
		Body      int       `json:"body"`
		Error     string    `json:"error"`
	}

	if len(params.Body) > 140 {
		data := returnVals{
			Error: "Chirp is too long",
		}
	}
}
