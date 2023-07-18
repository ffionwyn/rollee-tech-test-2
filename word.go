package main

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
	"sync"
)

var (
	storage     = make(map[string]int)
	storageLock sync.RWMutex
)

func main() {
	http.HandleFunc("/store", handleStoreWord)
	http.HandleFunc("/search", handleSearchWord)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleStoreWord(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	word := r.FormValue("word")
	if !isValidWord(word) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Incorrect format")
		return
	}

	word = strings.ToLower(word)
	storageLock.Lock()
	defer storageLock.Unlock()
	storage[word]++

	fmt.Fprintf(w, "'%s' stored successfully", word)
}

func handleSearchWord(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	prefix := r.FormValue("prefix")
	if !isValidWord(prefix) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Invalid prefix format")
		return
	}

	prefix = strings.ToLower(prefix)
	maxCount := 0

	storageLock.RLock()
	defer storageLock.RUnlock()
	for word, count := range storage {
		wordLower := strings.ToLower(word)
		if strings.HasPrefix(wordLower, prefix) && count > maxCount {
			maxCount = count
		}
	}

	if maxCount == 0 {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, "No matching word found in the storage")
		return
	}
	fmt.Fprintf(w, "Most frequent word with prefix '%s': %s", prefix, getMaxWordWithPrefix(prefix))
}

func getMaxWordWithPrefix(prefix string) string {
	var maxWord string
	maxCount := 0

	storageLock.RLock()
	defer storageLock.RUnlock()
	for word, count := range storage {
		wordLower := strings.ToLower(word)
		if strings.HasPrefix(wordLower, prefix) && count > maxCount {
			maxWord = word
			maxCount = count
		}
	}
	return maxWord
}

func isValidWord(word string) bool {
	validWordRegex := regexp.MustCompile(`^([a-zA-Z])+$`)
	return validWordRegex.MatchString(word)
}