package util

import (
	"io"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2/log"
)

func File_read_text(file_name string) ([]string, error) {

	// open file
	file, err := os.Open(file_name)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = file.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	// read file
	b, err := io.ReadAll(file)

	// done
	return strings.Split(string(b), "\n"), err
}
