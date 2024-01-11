package service

import (
	"bufio"
	"fmt"
	"os"
)

const (
	wordsBank      = "https://raw.githubusercontent.com/dwyl/english-words/master/words.txt"
	topWordsAmount = 10
)

// ReadURLsFromFile read all url's from a given filePath. we use go routines  for better performance
func ReadURLsFromFile(filePath string, numWorkers int) ([]string, error) {
	var urls []string
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		url := scanner.Text()
		urls = append(urls, url)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Reached end Of file", err)
	}

	return urls, nil
}