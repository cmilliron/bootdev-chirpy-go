package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	const (
		port = "8080"
		filePathRoot = "."
	)
	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir(filePathRoot))
	mux.Handle("/", fileServer)
	server := &http.Server{
		Addr: ":" + port,
		Handler: mux,
	}


	fmt.Printf("Server files from %s on port %s\n", filePathRoot, server.Addr)

	log.Fatal(server.ListenAndServe())

	defer server.Close()
}