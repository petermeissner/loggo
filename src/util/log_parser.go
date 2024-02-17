package util

import (
	"errors"
	"regexp"
)

type Log_message struct {
	Log_level string
	Message   string
}

type Log_file_parser struct {
	Line_input              func() (string, error)
	Re_start                *regexp.Regexp
	Max_n_messages          uint64
	N_messages              uint64
	Max_n_lines             uint
	Internal_status         string
	Internal_status_message string
	line_buffer             string
	message_buffer          string
	message                 string
}

func (log_file_parser *Log_file_parser) Status() (string, error) {
	if (log_file_parser.N_messages >= log_file_parser.Max_n_messages) || (log_file_parser.Internal_status == "done") {
		return "done", nil
	} else if log_file_parser.Internal_status == "error" {
		return "error", errors.New("Error")
	} else {
		return "ok", nil
	}
}

func (log_file_parser *Log_file_parser) Next() (string, error) {
	if log_file_parser.Internal_status != "ok" {
		return "", errors.New("Error")
	} else if (log_file_parser.message != "") && (log_file_parser.Internal_status == "ok") {
		return log_file_parser.message, nil
	} else {
		log_file_parser.Parse()
		return log_file_parser.Next()
	}
}

func (log_file_parser *Log_file_parser) Parse() {
	for {
		// get new line if buffer is empty or use buffer
		line, e := log_file_parser.Line_input()

		// handle error
		if e != nil {
			log_file_parser.Internal_status = "error"
			log_file_parser.Internal_status_message = e.Error()
			break
		}

		// new entry (starts with date)
		is_new_entry := log_file_parser.Re_start.Match([]byte(line))

		// handle state
		// 1. start new entry => message buffer is empty & current line matches start of new entry
		if is_new_entry && log_file_parser.message_buffer == "" {
			log_file_parser.message_buffer = line
			log_file_parser.line_buffer = ""
			continue
		} else
		// 2. finish old entry & start new entry => message buffer is not empty & current line matches start of new entry
		if is_new_entry && log_file_parser.message_buffer != "" {
			log_file_parser.message = log_file_parser.message_buffer + "\n" + line
			log_file_parser.line_buffer = ""
			log_file_parser.message_buffer = ""
			log_file_parser.N_messages++
			break
		} else
		// 3. continue old entry => message buffer is not empty & current line does not match start of new entry
		if !is_new_entry && log_file_parser.message_buffer != "" {
			log_file_parser.message_buffer += "\n" + line
			log_file_parser.line_buffer = ""
			continue
		} else
		// 4. continue old entry => message buffer is empty & current line does not match start of new entry
		if !is_new_entry && log_file_parser.message_buffer == "" {
			log_file_parser.line_buffer = line
			continue
		}
	}
}
