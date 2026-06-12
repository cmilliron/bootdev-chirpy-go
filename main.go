package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/cmilliron/bootdev-chirpy-go/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileServerHits	atomic.Int32
	db		*database.Queries
	platform	string
}

func main() {
	const (
		port = "8080"
		filePathRoot = "./static"
	)
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL must be set")
	}
	
	platform := os.Getenv("PLATFORM")
	if platform == "" {
		log.Fatal("PLATFORM must be set")
	}
	fmt.Println("url: " + dbURL)
	db, err := sql.Open("postgres", dbURL)
	if (err != nil) {
		log.Fatal("Database Connection Failed")
	}
	dbQueries := database.New(db)

	apiCfg := apiConfig{
		db: dbQueries,
		platform: platform,
	}
	
	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir(filePathRoot))
	
	// strip the app so that it doesn't automatically the route to the path and start serving 
	// from static/ and not try to serve from app/static/
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", fileServer)))

	// api routes
	mux.HandleFunc("GET /api/healthz", healthStatusHandler)
	// mux.HandleFunc("POST /api/validate_chirp", handleValidateChirp)
	mux.HandleFunc("POST /api/chirps", apiCfg.handlerCreateChirp)
	mux.HandleFunc("GET /api/chirps", apiCfg.handlerGetAllChrips)
	mux.HandleFunc("GET /api/chirps/{ChirpId}", apiCfg.handlerSingleChirp)
	mux.HandleFunc("POST /api/users", apiCfg.handlerCreateUser)
	mux.HandleFunc("POST /api/login", apiCfg.handlerLogin)

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
	if (cfg.platform != "dev") {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Reset is only allowed in dev environment."))
		return
	}
	
	cfg.fileServerHits.Store(0)
	
	err := cfg.db.Reset(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to reset the database: " + err.Error()))
		return
	}
	
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
