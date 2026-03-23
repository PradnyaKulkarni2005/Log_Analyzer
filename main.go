package main // entry point for the program
import (
	"bufio" // for reading files
	"fmt" // for printing output
	"os" // work with files - open, close
)

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
	for scanner.Scan(){
		line :=scanner.Text()
		fmt.Println(line)
		
	}


}