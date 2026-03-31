package handler

// handler package to handle HTTP requests for log analysis
import (
	"encoding/json"
	"io"
	"log-analyzer/internal/analyzer"
	"net/http"
)

// handles the /analyze endpoint, reads the request body, analyzes the logs, and returns the results.
// http.ResponseWriter is used to write the response back to the client, and *http.Request contains the incoming request data.
func AnalyzeHandler(w http.ResponseWriter, r *http.Request) {
	// Read the request body
	body, err := io.ReadAll(r.Body)
	//error handling
	if err != nil {
		http.Error(w, "Failed to read input", http.StatusInternalServerError)
		return
	}
	// get the results from analyzer package
	result := analyzer.AnalyzeWithPerformance(string(body))
	// write JSON response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
