package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"
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

	// read file line by line
	i := 0
	for scanner.Scan() {
		fmt.Println(scanner.Text())

		// wait for 1 second
		time.Sleep(250 * time.Millisecond)
		if i > 1000 {
			break
		}
		i++
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
