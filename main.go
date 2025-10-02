package main

import "net/http"

func main() {
	serveMux := http.NewServeMux()
	server := http.Server{}
	server.Handler = serveMux
	server.Addr = ".8080"

	server.ListenAndServe()
}
