package main

import (
	"regexp"

	"github.com/petermeissner/loggo/src/util"
)

func read_log_messages(file_name string, pattern string) ([]string, error) {

	// check that regex compiles
	regex_start := regexp.MustCompile(pattern)

	// read file
	text, e := util.File_read_text(file_name)
	if e != nil {
		return nil, e
	}

	// parse log messages
	s := parse_log_messages(text, regex_start)

	// done
	return s, nil
}
