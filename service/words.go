package service

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

// LoadValidWords reads all words from the URL and returns the words which are valid according to the requirements.
func LoadValidWords() map[string]struct{} {
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
	validWords := make(map[string]struct{})
	var mu sync.Mutex
	for _, word := range strings.Split(words, "\n") {
		if IsValidWord(word) {
			mu.Lock()
			validWords[word] = struct{}{}
			mu.Unlock()
		}
	}
	return validWords
}

// CountWordsFromAllURLs reads all words from all urls and put them in one global map. key=word value=count
// use sliding window rate limiter for sending http requests
func CountWordsFromAllURLs(urls []string, validWords map[string]struct{}) map[string]int {
	countWordsMap := make(map[string]int)
	var wg sync.WaitGroup
	var lock sync.Mutex
	for _, url := range urls {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()

			//wait until a request is allowed within the window
			for !allowRequest() {
				time.Sleep(10 * time.Millisecond) // Adjust sleep duration as needed
			}
			resp, err := http.Get(url)
			if err != nil {
				fmt.Println("Error fetching URL:", err)
			} else {
				defer resp.Body.Close()

				scanner := bufio.NewScanner(resp.Body)
				for scanner.Scan() {
					line := scanner.Text()
					// Split line into words (being generic)
					words := strings.Fields(line)

					for _, word := range words {
						if _, ok := validWords[word]; ok {
							lock.Lock()
							countWordsMap[word]++
							lock.Unlock()
						}
					}
				}
			}
		}(url)
	}

	wg.Wait()
	return countWordsMap

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
	minCount := int(^uint(0) >> 1)
	for word, count := range wordCounts {
		if count < minCount {
			minWord = word
			minCount = count
		}
	}
	return minWord, minCount
}
