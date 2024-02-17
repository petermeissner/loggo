package main

import (
	"regexp"
	"testing"
)

func TestParseLogMessages_single_message(t *testing.T) {

	// test data
	text := []string{`2021-08-01 12:00:00,000 [main] INFO  com.example.Main - Application started`}

	// parse
	messages := parse_log_messages(text, regexp.MustCompile(`^\d{4}-\d{2}-\d{2}`))

	// check
	if len(messages) != 1 {
		t.Error("unexpected number of messages")
	}
}

func TestParseLogMessages_multiple_messages(t *testing.T) {

	// test data
	text := []string{`2021-08-01 12:00:00,000 [main] INFO  com.example.Main - Application started`, `2021-08-01 12:00:00,000 [main] INFO  com.example.Main - Application started`, `2021-08-01 12:00:00,000 [main] INFO  com.example.Main - Application started`, `2021-08-01 12:00:00,000 [main] INFO  com.example.Main - Application started`}

	// parse
	messages := parse_log_messages(text, regexp.MustCompile(`^\d{4}-\d{2}-\d{2}`))

	// check
	if len(messages) != 4 {
		t.Error("unexpected number of messages")
	}
}

func TestParseLogMessages_broken(t *testing.T) {

	// test data
	text := []string{`broken message`, `2021-08-01 12:00:00,000 [main] INFO  com.example.Main - Application started`, `2021-08-01 12:00:00,000 [main] INFO  com.example.Main - Application started`, `2021-08-01 12:00:00,000 [main] INFO  com.example.Main - Application started`}

	// parse
	messages := parse_log_messages(text, regexp.MustCompile(`^\d{4}-\d{2}-\d{2}`))

	// check
	if len(messages) != 4 {
		t.Error("unexpected number of messages")
	}
}

func TestParseLogMessages_empty_messages(t *testing.T) {

	// test data
	text := []string{}

	// parse
	messages := parse_log_messages(text, regexp.MustCompile(`^\d{4}-\d{2}-\d{2}`))

	// check
	if len(messages) != 0 {
		t.Error("unexpected number of messages")
	}
}

func TestParseLogMessages_empty_message(t *testing.T) {

	// test data
	text := []string{""}

	// parse
	messages := parse_log_messages(text, regexp.MustCompile(`^\d{4}-\d{2}-\d{2}`))

	// check
	if len(messages) != 1 {
		t.Error("unexpected number of messages")
	}
}
