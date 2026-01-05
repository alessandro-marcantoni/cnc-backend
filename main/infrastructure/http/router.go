package http

import (
	"net/http"
)

func NewRouter() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/v1.0/health", HealthHandler)
	mux.HandleFunc("/api/v1.0/members", MembersHandler)
	mux.HandleFunc("/api/v1.0/members/", MemberByIDHandler)

	router := cors(mux)
	return withMiddleware(router)
}
