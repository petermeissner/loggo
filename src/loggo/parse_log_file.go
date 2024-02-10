package main

import (
	"bufio"
	"log"
	"os"
	"regexp"
)

func parse_log_file(f_name string) {

	// open file
	file, err := os.Open(f_name)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// file reader
	scanner := bufio.NewScanner(file)
	buf := make([]byte, 1e9)
	scanner.Buffer(buf, 1e9)
	regex_start_date, _ := regexp.Compile("^[0-9]{4}-[0-9]{2}-[0-9]{2}") // regex to match date
	regex_start_at, _ := regexp.Compile("^ +at ")                        // regex to match "at" in stack trace

	// read file line by line
	var start_date = false
	var log_entries = []string{}
	var log_entry = ""
	var start_at = 0

	for scanner.Scan() {
		// get line
		line := scanner.Text()

		// check lines
		start_date = regex_start_date.Match([]byte(line))
		if regex_start_at.Match([]byte(line)) {
			start_at = start_at + 1
		} else {
			start_at = 0
		}

		// determine how to handle line
		// - start new entry
		// - append to current entry
		// - discard
		if start_date {
			log_entries = append(log_entries, log_entry)
			log_entry = line
		} else {
			if start_at <= 1 {
				log_entry = log_entry + "\n" + line
			}
		}
	}

	println(log_entries)

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
