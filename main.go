package main

import (
	"fmt"
	"wordscount/service"
)

const urlsFile = "endg-urls"

func main() {
	urls, err := service.ReadURLsFromFile(urlsFile)
	if err != nil {
		fmt.Println("Error reading URLs:", err)
		return
	}
	fmt.Printf("There are %d URL's \n", len(urls))

	//Load valid words from the word bank (not all the words in the word bank are valid)
	validWords := service.LoadValidWords()
	fmt.Printf("There are %d  valid words\n", len(validWords))

	wordCountMap := service.CountWordsFromAllURLs(urls, validWords)

	// Find the top configured number of words
	topWords := service.FindTopFrequentWords(wordCountMap)

	// Print pretty JSON output
	service.PrintInJson(topWords)

}
