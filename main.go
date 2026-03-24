package main // entry point for the program
import (
	"bufio" // for reading files
	"fmt" // for printing output
	"os" // work with files - open, close
	"strings" // for string manipulation
)
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
		fmt.Println("Error opening file:", err) // print the error message
		return // exit the program
	}
	defer file.Close() // close the file when the function finishes
	scanner := bufio.NewScanner(file)
	// this will create a scanner object that will read the file line by line i.e doesnt load the entire file into memory at once
	errorCount := 0
	infoCount := 0
	ipCount := make(map[string]int) // create a map to count the occurrences of each IP address

	for scanner.Scan(){
		line :=scanner.Text()
		if isErrorLine(line) {
			errorCount++
		}
		if isInfoLine(line) {
			infoCount++
		}
		if ip,ok := extractIP(line); ok {
			ipCount[ip]++ // increment the count for this IP address
		}

	}
	fmt.Printf("Total number of errors: %d\n", errorCount) // print the total number of errors
	fmt.Printf("Total number of info messages: %d\n", infoCount) // print the total number of info messages

	fmt.Println("IP address counts:") // print the header for the IP address counts
	for ip,count :=range ipCount {
		fmt.Println(ip, ":", count) 
	}
}