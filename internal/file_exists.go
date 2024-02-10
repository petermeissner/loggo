package internal

import "os"

func File_exists(path string) bool {
	if _, err := os.Stat(path); err == nil {
		// exists
		return true
	}
	return false
}
