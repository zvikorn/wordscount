package service

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
)

// LoadValidWords reads all words from the URL and returns the words which are valid according to the requirements.
func LoadValidWords() map[string]bool {
	response, err := http.Get(wordsBank)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	words := string(body)
	validWords := make(map[string]bool)
	var mu sync.Mutex
	for _, word := range strings.Split(words, "\n") {
		if IsValidWord(word) {
			mu.Lock()
			validWords[word] = true
			mu.Unlock()
		}
	}
	return validWords
}

// CountWordsFromAllURLs reads all words from all urls and put them in one global map. key=word value=count
func CountWordsFromAllURLs(urls []string, validWords map[string]bool) map[string]int {
	// Create a channel for collecting word counts
	wordCountCh := make(chan map[string]int, 10)
	var wg sync.WaitGroup

	// Start a goroutine to wait for all workers to finish and close wordCountsCh
	go func() {
		wg.Wait()
		close(wordCountCh)
	}()

	//TODO instead spawning so many routines, consider having less routines, each will handle several urls
	for _, url := range urls {
		wg.Add(1)
		go countWords(url, validWords, wordCountCh, &wg)
	}

	//global word count map, holds all words and counts from all other maps
	countWordsMap := make(map[string]int)
	var mu sync.Mutex
	for wordCount := range wordCountCh {
		//we iterate each wordCount (from different routine) and update the global map
		mu.Lock()
		for k, v := range wordCount {
			countWordsMap[k] += v
		}
		mu.Unlock()
	}

	return countWordsMap

}

// countWords counts the valid words included in the url content, save them in a map channel (wordCounts).
// The key in the map is the word and the value is the count of this word.
// TODO use some kind of rate limiter
func countWords(url string, validWords map[string]bool, wordCountCh chan map[string]int, wg *sync.WaitGroup) {
	defer wg.Done()
	response, err := http.Get(url)
	if err != nil {
		fmt.Printf("Could not get url %s\n", url)
		return
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("Could not read response from url %s", url)
	}

	tempMap := make(map[string]int)
	for _, word := range strings.Fields(string(body)) {
		word = strings.ToLower(word)
		if IsValidWord(word) && validWords[word] {
			tempMap[word]++
		}
	}
	wordCountCh <- tempMap

	return
}

// findTopWordsFrequentWords returns top configured number of words from the combinedCounts map
func FindTopFrequentWords(combinedCounts map[string]int) map[string]int {
	topWords := make(map[string]int, topWordsAmount)
	for word, count := range combinedCounts {
		if len(topWords) < topWordsAmount {
			topWords[word] = count
		} else {
			minWord, minCount := findMin(topWords)
			if count > minCount {
				delete(topWords, minWord)
				topWords[word] = count
			}
		}
	}
	return topWords
}

// findMin iterates the map of words and their count and find the word, for which the count is the minimum
// Note: we need to think if we can do it in a better way (performance)
func findMin(wordCounts map[string]int) (string, int) {
	minWord := ""
	minCount := int(^uint(0) >> 1) // Initialize to max int
	for word, count := range wordCounts {
		if count < minCount {
			minWord = word
			minCount = count
		}
	}
	return minWord, minCount
}
