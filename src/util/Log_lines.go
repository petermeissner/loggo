package util

import (
	"log"
	"strings"
)

func Log_lines(prefix string, s string, logger *log.Logger) {
	for _, line := range strings.Split(s, "\n") {
		logger.Println(prefix + line)
	}
}
