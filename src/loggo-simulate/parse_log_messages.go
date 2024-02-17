package main

import "regexp"

func parse_log_messages(text []string, regex_start *regexp.Regexp) []string {

	// go through lines and group lines into messages
	messages := []string{}
	for _, line := range text {
		if regex_start.MatchString(line) || len(messages) == 0 {
			messages = append(messages, line)
		} else {
			messages[len(messages)-1] = messages[len(messages)-1] + "\n" + line
		}
	}

	// done
	return messages
}
