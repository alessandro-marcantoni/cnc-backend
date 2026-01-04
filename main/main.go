package main

import (
	"log"
	"net/http"
	"time"

	internalHttp "github.com/alessandro-marcantoni/cnc-backend/main/infrastructure/http"
	"github.com/alessandro-marcantoni/cnc-backend/main/infrastructure/persistence"
)

func main() {
	// Initialize database connection
	db, err := persistence.InitializeDatabase()
	if err != nil {
		log.Fatalf("❌ Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Initialize services with database
	internalHttp.InitializeServices(db)

	mux := internalHttp.NewRouter()

	server := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Println("🚀 API listening on :8080")
	log.Fatal(server.ListenAndServe())
}
