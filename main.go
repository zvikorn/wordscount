package main

import (
	"fmt"
	"wordscount/service"
)

const urlsFile = "endg-urls"
const numWorkers = 4

func main() {
	urls, err := service.ReadURLsFromFile(urlsFile, numWorkers)
	if err != nil {
		fmt.Println("Error reading URLs:", err)
		return
	}

	//Load valid words from the word bank (not all the words in the word bank are valid)
	validWords := service.LoadValidWords()

	wordCountMap := service.CountWordsFromAllURLs(urls, validWords)

	// Find the top configured number of words
	topWords := service.FindTopFrequentWords(wordCountMap)

	// Print pretty JSON output
	service.PrintInJson(topWords)

}
