package analyzer

type Result struct {
	ErrorCount int
	InfoCount  int
	IP         string
	HasIP      bool
}

type FinalResult struct {
	TotalErrors int
	TotalInfo   int
	IPCount     map[string]int
}