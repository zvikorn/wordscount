package service

import (
	"encoding/json"
	"fmt"
	"regexp"
)

// IsValidWord given a word, checks if it's considered to be a valid word according to specified regex
func IsValidWord(word string) bool {
	// Ensure at least 3 characters and only alphabetic characters
	// NOTE: if the regex is an overhead we can check each char with ascii value (< >).
	match, err := regexp.MatchString("^[a-zA-Z]{3,}$", word)
	if err != nil {
		fmt.Printf("received error while validating the word %s: %v", word, err)
	}
	return match
}

// convert the received top x words to json and ptine them
func PrintInJson(topWords map[string]int) {
	jsonData, err := json.MarshalIndent(topWords, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(jsonData))
}
