package util

import (
	"io/ioutil"
	"log"
	"regexp"
)

func File_list(path string, pattern string) []string {
	// list files
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	// compile regex
	regex, err := regexp.Compile(pattern)
	if err != nil {
		log.Fatal(err)
	}

	// filter files
	var file_list = []string{}

	for _, file := range files {
		if !file.IsDir() && regex.Match([]byte(file.Name())) {
			file_list = append(file_list, path+"/"+file.Name())
		}
	}

	// return
	return file_list
}
