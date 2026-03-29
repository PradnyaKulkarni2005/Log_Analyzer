package analyzer
// analyzer package to analyze log data and extract insights
import (
	"strings"
	"sync"
	"log-analyzer/utils"
)
// worker function processes log lines from the jobs channel, analyzes them, and sends results to the results channel. It uses a WaitGroup to signal when it's done processing.

func worker(jobs <-chan string, results chan<- Result, wg *sync.WaitGroup) {
	// Ensure that the WaitGroup counter is decremented when the worker finishes processing
	defer wg.Done()
// Continuously read from the jobs channel until it's closed
	for line := range jobs {
		res := Result{}

		if utils.IsErrorLine(line) {
			res.ErrorCount = 1
		}

		if utils.IsInfoLine(line) {
			res.InfoCount = 1
		}

		if ip, ok := utils.ExtractIP(line); ok {
			res.IP = ip
			res.HasIP = true
		}

		results <- res
	}
}

func AnalyzeLogs(data string) FinalResult {
	// Split the input data into lines for processing
	lines := strings.Split(data, "\n")
	// Create channels for jobs and results, with a buffer size of 100 to allow for concurrent processing
	jobs := make(chan string, 100)
	results := make(chan Result, 100)
// Use a WaitGroup to manage the lifecycle of worker goroutines
	var wg sync.WaitGroup
	numWorkers := 3
// Start a fixed number of worker goroutines to process the log lines concurrently
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(jobs, results, &wg)
	}
// Send log lines to the jobs channel for processing by the workers, and close the channel once all lines have been sent
	go func() {
		for _, line := range lines {
			jobs <- line
		}
		close(jobs)
	}()
// Wait for all workers to finish processing and then close the results channel to signal that no more results will be sent
	go func() {
		wg.Wait()
		close(results)
	}()
// Aggregate results from the results channel - total error and info counts, IP address counts
	final := FinalResult{
		IPCount: make(map[string]int),
	}
// Read results from the results channel and update the final result accordingly
	for res := range results {
		final.TotalErrors += res.ErrorCount
		final.TotalInfo += res.InfoCount

		if res.HasIP {
			final.IPCount[res.IP]++
		}
	}

	return final
}