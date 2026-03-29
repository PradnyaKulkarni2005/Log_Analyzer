package main
// main entry point 
import (
	"fmt"
	"net/http"

	"log-analyzer/internal/handler"
)

func main() {
	// Register the handler for the /analyze endpoint
	// HandleFunc registers the handler function for the given pattern in the DefaultServeMux.
	http.HandleFunc("/analyze", handler.AnalyzeHandler)

	fmt.Println("Server running on :9090")
// Start the server on port 9090
// ListenAndServe listens on the TCP network address and then calls Serve with handler to handle requests on incoming connections.
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		fmt.Println("Server failed:", err)
	}
}