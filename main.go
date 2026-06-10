package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

func main() {
	const (
		port = "8080"
		filePathRoot = "./static"
	)
	apiCfg := apiConfig{}
	
	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir(filePathRoot))
	
	// strip the app so that it doesn't automatically the route to the path and start serving 
	// from static/ and not try to serve from app/static/
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", fileServer)))
	mux.HandleFunc("GET /healthz", healthStatusHandler)
	mux.HandleFunc("GET /metrics", apiCfg.metricsHandler)
	mux.HandleFunc("POST /reset", apiCfg.resetHandler)
	
	server := &http.Server{
		Addr: ":" + port,
		Handler: mux,
	}

	fmt.Printf("Server files from %s on port %s\n", filePathRoot, server.Addr)
	log.Fatal(server.ListenAndServe())
	defer server.Close()
}

type apiConfig struct {
	fileServerHits	 atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileServerHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) metricsHandler(w http.ResponseWriter, r *http.Request) {
	// printRequestHeader(r)
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	bodyContent := fmt.Sprintf("Hits: %d\n", cfg.fileServerHits.Load())
	body := []byte(bodyContent)
	w.Write(body)

}
func (cfg *apiConfig) resetHandler(w http.ResponseWriter, r *http.Request) {
	cfg.fileServerHits.Store(0)
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	bodyContent := fmt.Sprintf("Hits reset to 0.\nHits: %d\n", cfg.fileServerHits.Load())
	body := []byte(bodyContent)
	w.Write(body)
}


func healthStatusHandler(w http.ResponseWriter, r *http.Request) {
	// printRequestHeader(r)

	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	body := []byte(http.StatusText(http.StatusOK))
	w.Write(body)
}




// func printRequestHeader(r *http.Request) error {
// 	if r == nil {
// 		return fmt.Errorf("the request was malformed.\n")
// 	}
// 	for key, value := range r.Header {
// 		fmt.Printf("%s: %s\n", key, value)
// 	}
// 	return nil
// }

func middlewareLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s, %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}