package analyzer

type Result struct {
	ErrorCount int
	InfoCount  int
	IP         string
	HasIP      bool
	SlowRequest string
	LogsPerMinute map[string]int
}

type FinalResult struct {
	TotalErrors int
	TotalInfo   int
	TotalLogs   int
	ErrorRate float64
	IPCount map[string]int
	TopIP   string
	SuspiciousIPs []string
	SlowRequests []string
	LogsPerMinute map[string]int
	Performance Performance
}

type Performance struct {
	SequentialTimeMs  int64
	ConcurrentTimeMs  int64
	Speedup           float64
}