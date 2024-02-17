package main

import (
	"slices"
	"time"

	"github.com/petermeissner/loggo/src/util"
)

func max_mpm_handler(max_mpm float64, max_mpm_is_set bool, ts_pattern string, messages []string) float64 {

	if !max_mpm_is_set {
		// calculate messages per minute rate
		max_mpm = mpm_in_messages(ts_pattern, messages)

		// if rate is smaller than 60, set to 60
		if max_mpm < 60 {
			max_mpm = 60
		}
	}

	// done
	return max_mpm
}

func mpm_in_messages(ts_pattern string, messages []string) float64 {
	// variables
	var (
		tp              = util.TimePattern(ts_pattern)
		tp_layout       = tp.GolangTimeLayout()
		tp_re           = tp.ToRegex()
		mpm_values_unix []int64
		max_mpm         = float64(60)
	)

	// collect time stamps
	for _, message := range messages {
		msg_time, e := time.Parse(tp_layout, tp_re.FindString(message))
		if e == nil {
			mpm_values_unix = append(mpm_values_unix, msg_time.Unix())
		}
	}

	// calculate messages per minute rate
	if len(mpm_values_unix) >= 2 {
		max_mpm = float64(len(messages)) / (float64(slices.Max(mpm_values_unix)-slices.Min(mpm_values_unix)) / 60)
	}

	// done
	return max_mpm
}
