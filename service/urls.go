package service

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"sync"
)

const (
	wordsBank      = "https://raw.githubusercontent.com/dwyl/english-words/master/words.txt"
	topWordsAmount = 10
	numRoutines    = 4
)

// ReadURLsFromFile read all url's from a given filePath. we use go routines to have them read different
// segments of the file
func ReadURLsFromFile(filePath string) ([]string, error) {
	var urls []string

	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	fileSize, err := file.Stat()
	if err != nil {
		panic(err)
	}
	sectionSize := fileSize.Size() / int64(numRoutines)

	// these lock and counter are used just to show that the entire file is read
	var lock sync.Mutex
	// number of urls
	var count int64

	var wg sync.WaitGroup
	for i := 0; i < numRoutines; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			start := int64(i) * sectionSize
			end := start + sectionSize

			// Buffer for the segment
			buf := make([]byte, end-start)
			// Read the segment bytes
			_, err := file.ReadAt(buf, start)
			if err != nil && err != io.EOF {
				fmt.Println("Error reading segment:", err)
				return
			}

			// Create a reader from the bytes
			reader := bytes.NewReader(buf)
			scanner := bufio.NewScanner(reader)
			for scanner.Scan() {
				lock.Lock()
				url := scanner.Text()
				count++
				urls = append(urls, url)
				lock.Unlock()
			}
		}(i)
	}

	wg.Wait()
	return urls, nil
}
