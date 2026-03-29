package utils
// utils package to provide helper functions for log analysis
import "strings"

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