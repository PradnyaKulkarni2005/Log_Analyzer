package main // entry point for the program
import (
	"bufio" // for reading files
	"fmt" // for printing output
	"os" // work with files - open, close
	"strings" // for string manipulation
	"sync" // for synchronization primitives like WaitGroup
)

type Result struct {
	errorCount int
	infoCount int
	ip string
	hasIP bool
}
// jobs channnel - read only channel to send log lines to worker goroutines
// results channel - write only channel to receive counts from worker goroutines
// worker - receives line , processes it and sends results back
func worker(jobs <-chan string , results chan<- Result , wg *sync.WaitGroup){
	defer wg.Done() // signal that the worker is done when the function finishes
	// when this worker finishes , reduce counter by 1
	for line := range jobs {
		res := Result{}
		if isErrorLine(line) {
			res.errorCount=1
		}
	    if isInfoLine(line){
			res.infoCount=1
		}
		if ip , ok:= extractIP(line); ok{
			res.ip=ip
			res.hasIP=true
		}
		results <-res
}
}

func isErrorLine(line string) bool {
	return strings.Contains(line, "ERROR") // check if the line contains the word "ERROR"
}
func isInfoLine(line string) bool {
	return strings.Contains(line, "INFO")
}

func extractIP(line string) (string,bool){
	if strings.Contains(line, "IP="){
		parts:=strings.Split(line, "IP=") // split the line into parts using "IP=" as the delimiter
		return parts[1], true // return the second part of the split (the IP address) and true to indicate that an IP was found
	}
	return "", false // return an empty string and false to indicate that no IP was found
}
func main() {
	file, err := os.Open("log.txt") // open the file
	// returns a file object and an error object
	// go doesnt use exceptions, instead it returns an error object that we can check
	// if error comes then error handling

	if err != nil {
		fmt.Println("Error opening file:", err)
		return 
	}
	defer file.Close() // close the file when the function finishes
	scanner := bufio.NewScanner(file)
	// this will create a scanner object that will read the file line by line i.e doesnt load the entire file into memory at once
	errorCount := 0
	infoCount := 0
	ipCount := make(map[string]int) // create a map to count the occurrences of each IP address
	
	jobs := make(chan string) // create a channel to send lines to worker goroutines , it will carry log lines

	results := make(chan Result) // create a channel to receive counts from worker goroutines 
	numWorkers :=3
	var wg sync.WaitGroup

	for i:=0; i<numWorkers; i++{
		wg.Add(1) // increment the WaitGroup counter for each worker
		go worker(jobs, results, &wg)//starts a gouroutine that runs the worker function, passing the jobs and results channels as arguments and a pointer to the WaitGroup
	}
	go func() {
	for scanner.Scan() {
		jobs <- scanner.Text() // send each line to workers
	}
	close(jobs) // tell workers: no more data
}()

	go func(){
		wg.Wait() // wait for all workers to finish
		close(results) // close the results channel after all workers are done
	}()
		
	for result := range results{
		errorCount += result.errorCount
		infoCount += result.infoCount
		if result.hasIP {
			ipCount[result.ip]++ // increment the count for this IP address in the map
		} // add the count from the worker to the total error count
	}
	fmt.Printf("Total number of errors: %d\n", errorCount) // print the total number of errors
	fmt.Printf("Total number of info messages: %d\n", infoCount) // print the total number of info messages

	fmt.Println("IP address counts:") // print the header for the IP address counts
	for ip,count :=range ipCount {
		fmt.Println(ip, ":", count) 
	}
}