package main

import (
	"fmt"
	"net/http"

	"log-analyzer/internal/handler"
)

// CORS middleware
func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Allow all origins (for development)
		w.Header().Set("Access-Control-Allow-Origin", "*")

		// Allowed methods
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")

		// Allowed headers
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// Handle preflight request
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {

	// Register route with CORS middleware
	http.Handle("/analyze", enableCORS(http.HandlerFunc(handler.AnalyzeHandler)))

	fmt.Println("🚀 Server running on http://localhost:9090")

	// Start server
	if err := http.ListenAndServe(":9090", nil); err != nil {
		fmt.Println("Server failed:", err)
	}
}