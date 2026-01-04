package main

import (
	"log"
	"net/http"
	"time"

	internalHttp "github.com/alessandro-marcantoni/cnc-backend/main/infrastructure/http"
)

func main() {
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
