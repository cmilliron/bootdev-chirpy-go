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

	// api routes
	mux.HandleFunc("GET /api/healthz", healthStatusHandler)
	mux.HandleFunc("POST /api/validate_chirp", handleValidateChirp)

	// admin routes
	mux.HandleFunc("GET /admin/metrics", apiCfg.metricsHandler)
	mux.HandleFunc("POST /admin/reset", apiCfg.resetHandler)
	
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
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	bodyContent := fmt.Sprintf(`<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`, cfg.fileServerHits.Load())
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
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	body := []byte(http.StatusText(http.StatusOK))
	w.Write(body)
}

func middlewareLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s, %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
