package utils
// utils package to provide helper functions for log analysis
import "strings"
import "strconv"

func IsErrorLine(line string) bool {
	return strings.Contains(line, "ERROR")
}

func IsInfoLine(line string) bool {
	return strings.Contains(line, "INFO")
}

func ExtractIP(line string) (string, bool) {
	if strings.Contains(line, "IP=") {
		parts := strings.Split(line, "IP=")
		ipPart := parts[1]
		ip := strings.Fields(ipPart)[0]
		return ip, true
	}
	return "", false
}

func ExtractResponseTime(line string) (int, bool) {
	// checks if the line contains "took" and extracts the response time in milliseconds
	if strings.Contains(line, "took") {
		// Split the line by "took" and trim whitespace to get the time part
		parts := strings.Split(line, "took")
		timePart := strings.TrimSpace(parts[1])
// Remove the "ms" suffix and convert the remaining string to an integer
		timeStr := strings.TrimSuffix(timePart, "ms")
		// Convert the time string to an integer value
		timeVal, err := strconv.Atoi(timeStr)
		if err == nil {
			return timeVal, true
		}
	}
	return 0, false
}


