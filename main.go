package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	const (
		port = "8080"
		filePathRoot = "./static"
	)
	
	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir(filePathRoot))
	
	// strip the app so that it doesn't automatically the route to the path and start serving 
	// from static/ and not try to serve from app/static/
	mux.Handle("/app/", http.StripPrefix("/app", fileServer))
	mux.HandleFunc("/healthz", healthStatusHandler)
	
	server := &http.Server{
		Addr: ":" + port,
		Handler: mux,
	}

	fmt.Printf("Server files from %s on port %s\n", filePathRoot, server.Addr)

	log.Fatal(server.ListenAndServe())

	defer server.Close()
}

func healthStatusHandler(w http.ResponseWriter, r *http.Request) {
	printRequestHeader(r)

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	body := []byte("OK")
	w.Write(body)
}

func printRequestHeader(r *http.Request) error {
	if r == nil {
		return fmt.Errorf("the request was malformed.\n")
	}
	for key, value := range r.Header {
		fmt.Printf("%s: %s\n", key, value)
	}
	return nil
}