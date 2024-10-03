package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/dogeorg/reflector/pkg/api"
	"github.com/dogeorg/reflector/pkg/database"
	reflectormiddleware "github.com/dogeorg/reflector/pkg/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func CORS() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "*")
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}

func main() {
	// Initialize database
	db, err := database.NewDatabase()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Create router
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
		reflectormiddleware.RateLimiter(time.Minute, 1)(api.CreateEntry(db)).ServeHTTP(w, r)
	})

	r.With(CORS()).Get("/dbxcheck", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("{\"status\": \"ok\"}"))
	})

	r.With(CORS()).Options("/{token}", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	r.With(CORS()).Get("/{token}", func(w http.ResponseWriter, r *http.Request) {
		api.GetIP(db).ServeHTTP(w, r)
	})

	bindAddr := flag.String("bind", ":8080", "Bind address and port for the server")
	flag.Parse()

	// Start server
	log.Printf("Server starting on %s\n", *bindAddr)
	log.Fatal(http.ListenAndServe(*bindAddr, r))
}
