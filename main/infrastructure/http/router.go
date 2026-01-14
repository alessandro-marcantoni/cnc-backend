package http

import (
	"net/http"
)

func NewRouter() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/v1.0/health", HealthHandler)
	mux.HandleFunc("/api/v1.0/members", MembersHandler)
	mux.HandleFunc("/api/v1.0/members/", MemberByIDHandler)
	mux.HandleFunc("/api/v1.0/memberships", MembershipsHandler)
	mux.HandleFunc("/api/v1.0/facilities/catalog", FacilitiesCatalogHandler)
	mux.HandleFunc("/api/v1.0/facilities", FacilitiesByTypeHandler)
	mux.HandleFunc("/api/v1.0/facilities/rented", RentedFacilitiesHandler)

	router := cors(mux)
	return withMiddleware(router)
}
