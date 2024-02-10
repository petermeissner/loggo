package main

import "os"

func file_exists(path string) bool {
	if _, err := os.Stat(path); err == nil {
		// exists
		return true
	}
	return false
}
