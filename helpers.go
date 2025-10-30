package main

import "net/http"

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

	}
}
