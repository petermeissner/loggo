package main

import (
	"bufio"
	"log"
	"os"
	"regexp"
	"strings"
)

func main() {

	// open file
	file, err := os.Open(".\\logs\\error.log")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// file reader
	scanner := bufio.NewScanner(file)
	buf := make([]byte, 1e9)
	scanner.Buffer(buf, 1e9)
	regex_start_date, _ := regexp.Compile("^[0-9]{4}-[0-9]{2}-[0-9]{2}") // regex to match date
	regex_start_at, _ := regexp.Compile(`^\s+at `)                       // regex to match "at" in stack trace
	regex_start_more, _ := regexp.Compile(`^\s+... \d+ more`)            // regex to match "... X more" in stack trace
	regex_start_caused, _ := regexp.Compile(`^\s*[Cc]aused by: `)
	regex_all_whitespace, _ := regexp.Compile(`^\s*$`)

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
			println(log_entry)
			log_entry = line
		} else {
			if start_at <= 1 && !regex_start_more.Match([]byte(line)) && !regex_all_whitespace.Match([]byte(line)) {
				if start_at == 1 {
					log_entry = log_entry + " " + strings.Replace(line, "\t", "", -1)
				} else {
					if regex_start_caused.Match([]byte(line)) {
						log_entry = log_entry + "\n  " + line
					} else {
						log_entry = log_entry + "\n" + line
					}
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
