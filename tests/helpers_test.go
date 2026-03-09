package tests

import "os"

// joinStrings concatenates strings with newlines
func joinStrings(strs []string) string {
	result := ""
	for _, str := range strs {
		result += str + "\n"
	}
	return result
}

// fileExists checks if a file exists at the given path
func fileExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// dirExists checks if a directory exists at the given path
func dirExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

// getOutputContent extracts string content from an Output interface
func getOutputContent(output interface{ Data() any }) string {
	data := output.Data()
	if bytes, ok := data.([]byte); ok {
		return string(bytes)
	}
	return ""
}
