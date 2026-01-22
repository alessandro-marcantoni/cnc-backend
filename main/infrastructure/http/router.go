package http

import (
	"net/http"
	"strings"
)

func NewRouter() http.Handler {
	mux := http.NewServeMux()

	// Health check - no auth required
	mux.HandleFunc("/api/v1.0/health", HealthHandler)

	// Protected routes - auth required
	mux.HandleFunc("/api/v1.0/members", MembersHandler)
	mux.HandleFunc("/api/v1.0/members/", MemberByIDHandler)
	mux.HandleFunc("/api/v1.0/memberships", MembershipsHandler)
	mux.HandleFunc("/api/v1.0/facilities/catalog", FacilitiesCatalogHandler)
	mux.HandleFunc("/api/v1.0/facilities", FacilitiesByTypeHandler)
	mux.HandleFunc("/api/v1.0/facilities/rented/", RentedFacilityByIDHandler)
	mux.HandleFunc("/api/v1.0/facilities/rented", RentedFacilitiesHandler)
	mux.HandleFunc("/api/v1.0/facilities/waiting-list", WaitingListHandler)
	mux.HandleFunc("/api/v1.0/facilities/suggested-price", SuggestedPriceHandler)
	mux.HandleFunc("/api/v1.0/payments", PaymentsHandler)
	mux.HandleFunc("/api/v1.0/payments/", PaymentByIDHandler)

	router := cors(mux)
	router = conditionalAuthMiddleware(router)
	router = loggingMiddleware(router)
	return withMiddleware(router)
}

// conditionalAuthMiddleware applies auth middleware to all routes except health
func conditionalAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip auth for health endpoint
		if strings.HasPrefix(r.URL.Path, "/api/v1.0/health") {
			next.ServeHTTP(w, r)
			return
		}

		if r.Method == http.MethodOptions {
			next.ServeHTTP(w, r)
			return
		}

		// Apply auth middleware for all other routes
		authMiddleware(next).ServeHTTP(w, r)
	})
}
