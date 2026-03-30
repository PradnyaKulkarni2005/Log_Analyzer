package analyzer
// analyzer package to analyze log data and extract insights
import (
	"strings"
	"sync"
	"log-analyzer/utils"
	"time"

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
		// slow request detection - set a threshold time and check if any log line contains a request that exceeds this threshold
		if time, ok := utils.ExtractResponseTime(line); ok {
	if time > 100 {
		res.SlowRequest = line
	}
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
		// add slow request to the final result if it exists
		if res.SlowRequest != "" {
	final.SlowRequests = append(final.SlowRequests, res.SlowRequest)
		}
		
	}

	// calculating error rate - (total errors / total logs ) * 100
	final.TotalLogs = len(lines)
	if final.TotalLogs > 0 {
	final.ErrorRate = float64(final.TotalErrors) / float64(final.TotalLogs) * 100
	}

	// calculating the top IP - most frequent IP address
	maxCount := 0
	for ip, count := range final.IPCount {
		if count > maxCount {
			maxCount = count
			final.TopIP = ip
		}
	}

	// suspicious IP detection
	//logic - set a threshold , if count of a particular IP exceeds the threshold, consider it suspicious
	threshold := 4
	for ip, count := range final.IPCount {
	if count > threshold {
		final.SuspiciousIPs = append(final.SuspiciousIPs, ip)
	}
	}

	// slow request detection - set a threshold time and check if any log line contains a request that exceeds this threshold
	


	return final
}

func AnalyzeLogsSequential(data string) FinalResult {
	lines := strings.Split(strings.TrimSpace(data), "\n")

	final := FinalResult{
		IPCount: make(map[string]int),
	}

	for _, line := range lines {

		if utils.IsErrorLine(line) {
			final.TotalErrors++
		}

		if utils.IsInfoLine(line) {
			final.TotalInfo++
		}

		if ip, ok := utils.ExtractIP(line); ok {
			final.IPCount[ip]++
		}

		// slow request
		if time, ok := utils.ExtractResponseTime(line); ok {
			if time > 100 {
				final.SlowRequests = append(final.SlowRequests, line)
			}
		}
	}

	final.TotalLogs = len(lines)

	return final
}
func AnalyzeWithPerformance(data string) FinalResult {

	// ⏱️ Sequential timing
	startSeq := time.Now()
	_ = AnalyzeLogsSequential(data) // ignore result
	seqTime := time.Since(startSeq).Milliseconds()

	// ⏱️ Concurrent timing
	startCon := time.Now()
	conResult := AnalyzeLogs(data)
	conTime := time.Since(startCon).Milliseconds()

	// 📊 Add performance data
	conResult.Performance = Performance{
		SequentialTimeMs: seqTime,
		ConcurrentTimeMs: conTime,
	}

	if conTime > 0 {
		conResult.Performance.Speedup = float64(seqTime) / float64(conTime)
	}

	return conResult
}