package main

import (
	"log"
	"net/http"
)

func (cfg *apiConfig) metricInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileServerHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func returnJSON(w http.ResponseWriter, code int, payload interface{}) {

}

func returnError(w http.ResponseWriter, code int, msg string) {
	if code > 499 {
		log.Printf("5xx error: %s", msg)
	}
	type errorReturn struct {
		Error string `json"error"`
	}
	returnJSON(w, code, errorReturn{
		Error: msg,
	})
}
