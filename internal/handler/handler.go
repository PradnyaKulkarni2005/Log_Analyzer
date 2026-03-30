package handler
// handler package to handle HTTP requests for log analysis
import (
	"fmt"
	"io"
	"net/http"
	"log-analyzer/internal/analyzer"
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
	result := analyzer.AnalyzeLogs(string(body))
	// write the results back to the client
	fmt.Fprintf(w, "Errors: %d\n", result.TotalErrors)
fmt.Fprintf(w, "Info: %d\n", result.TotalInfo)
fmt.Fprintf(w, "Total Logs: %d\n", result.TotalLogs)

fmt.Fprintf(w, "Error Rate: %.2f%%\n", result.ErrorRate)

fmt.Fprintf(w, "Top IP: %s\n", result.TopIP)

// Suspicious IPs
fmt.Fprintf(w, "Suspicious IPs:\n")
for _, ip := range result.SuspiciousIPs {
	fmt.Fprintf(w, "- %s\n", ip)
}

// Slow Requests
fmt.Fprintf(w, "Slow Requests:\n")
for _, req := range result.SlowRequests {
	fmt.Fprintf(w, "- %s\n", req)
}

// IP Counts (keep this at end)
fmt.Fprintf(w, "IP Counts:\n")
for ip, count := range result.IPCount {
	fmt.Fprintf(w, "%s : %d\n", ip, count)
}
}