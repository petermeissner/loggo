package util

import (
	"regexp"
	"strings"
	"time"
)

type TimePattern string

func (t TimePattern) String() string {
	return string(t)
}

func (t TimePattern) GolangTimeLayout() string {
	s := strings.ReplaceAll(string(t), "yyyy", "2006")
	s = strings.ReplaceAll(string(s), "yy", "06")
	s = strings.ReplaceAll(string(s), "MM", "01")
	s = strings.ReplaceAll(string(s), "dd", "02")
	s = strings.ReplaceAll(string(s), "HH", "15")
	s = strings.ReplaceAll(string(s), "mm", "04")
	s = strings.ReplaceAll(string(s), "ss", "05")
	s = strings.ReplaceAll(string(s), "S", "9")
	return s
}

func (t TimePattern) ToRegex() regexp.Regexp {
	s := strings.ReplaceAll(string(t), "yyyy", "\\d{4}")
	s = strings.ReplaceAll(string(s), "yy", "\\d{2}")
	s = strings.ReplaceAll(string(s), "MM", "\\d{2}")
	s = strings.ReplaceAll(string(s), "dd", "\\d{2}")
	s = strings.ReplaceAll(string(s), "HH", "\\d{2}")
	s = strings.ReplaceAll(string(s), "mm", "\\d{2}")
	s = strings.ReplaceAll(string(s), "ss", "\\d{2}")
	s = strings.Replace(string(s), "S", "\\d+", 1)
	s = strings.ReplaceAll(string(s), "S", "")
	return *regexp.MustCompile(s)
}

func (t TimePattern) FormatTimeAsStringFunc() func(time.Time) string {
	return func(time time.Time) string {
		return time.Format(t.GolangTimeLayout())
	}
}
