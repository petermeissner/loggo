package main

import (
	"log"
	"os"
)

// return list of log available in certain directory
func list_log_files() []string {

	// list of log files
	var log_files []string

	// go through directory and get list of files ending with .log
	files, err := os.ReadDir("./logs/")
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if !file.IsDir() {
			log_files = append(log_files, file.Name())
		}
	}
	return log_files
}
