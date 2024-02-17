package main

import (
	"time"

	"github.com/petermeissner/loggo/src/util"
)

func update_time_stamp(messages []string, time_pattern string) []string {

	// setup time pattern
	var (
		tp             = util.TimePattern(time_pattern)
		tp_layout      = tp.GolangTimeLayout()
		tp_re          = tp.ToRegex()
		time_formatter = tp.FormatTimeAsStringFunc()
	)

	// slice to store messages
	m := []string{}

	// store first time found in messages
	var min_time time.Time

	// store start time
	start_time := time.Now()

	// go through messages and update timestamp
	for _, message := range messages {

		// find time stamp
		// - if not found just append message
		// - if found update time stamp (current time + offset within log time) and append message

		// try parsing time stamp
		msg_time, e := time.Parse(tp_layout, tp_re.FindString(message))

		// handle success and error
		if e == nil {

			// if timestamp found in message for the first time, store as min time for offset calculation
			if min_time.IsZero() {
				min_time = msg_time
			}

			// calculate new timestamp
			new_time := start_time.Add(msg_time.Sub(min_time))

			// replace old timestamp with new timestamp
			m = append(m, string(tp_re.ReplaceAll([]byte(message), []byte(time_formatter(new_time)))))

		} else {

			// error: leave message as is
			m = append(m, message)

		}
	}

	// done
	return m
}
