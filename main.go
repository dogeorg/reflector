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

	// Routes
	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
		reflectormiddleware.RateLimiter(time.Minute, 1)(api.CreateEntry(db)).ServeHTTP(w, r)
	})

	r.Get("/{token}", api.GetIP(db))

	// Parse command-line arguments
	bindAddr := flag.String("bind", ":8080", "Bind address and port for the server")
	flag.Parse()

	// Start server
	log.Printf("Server starting on %s\n", *bindAddr)
	log.Fatal(http.ListenAndServe(*bindAddr, r))
}
