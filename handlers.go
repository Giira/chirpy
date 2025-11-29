package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
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
	if cfg.platform != "dev" {
		writer.WriteHeader(403)
		return
	}
	cfg.queries.Reset(context.Background())
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
		msg := fmt.Sprintf("Error decoding parameters: %v", err)
		returnError(writer, 500, msg)
		return
	}

	params.Body = checkProfanity(params.Body)

	type returnVals struct {
		CleanedBody string `json:"cleaned_body"`
	}

	if len(params.Body) > 140 {
		msg := "Chirp is too long"
		returnError(writer, 400, msg)
	} else {
		data := returnVals{
			CleanedBody: params.Body,
		}
		returnJSON(writer, 200, data)
	}
}

func (cfg *apiConfig) handleCreateUser(writer http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Email string `json:"email"`
	}

	type user struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email     string    `json:"email"`
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		writer.WriteHeader(500)
		return
	}

	new_user, err := cfg.queries.CreateUser(context.Background(), params.Email)
	if err != nil {
		log.Printf("error creating user: %v", err)
		writer.WriteHeader(500)
		return
	}
	n_user := user{
		ID:        new_user.ID,
		CreatedAt: new_user.CreatedAt,
		UpdatedAt: new_user.UpdatedAt,
		Email:     new_user.Email,
	}

	returnJSON(writer, 201, n_user)
}
