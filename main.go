package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {

	port := "8080"
	root := "."

	serveMux := http.NewServeMux()
	file_handler := http.FileServer(http.Dir(root))
	serveMux.Handle("/app/", http.StripPrefix("/app", file_handler))

	serveMux.HandleFunc("/healthz", handleReady)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: serveMux,
	}

	fmt.Printf("Server active on port %v", port)
	log.Fatal(server.ListenAndServe())

}
