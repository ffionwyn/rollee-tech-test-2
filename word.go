package main

import (
	"fmt"
	"log"
	"net/http"
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

func handleService(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		handleStoreWord(w, r)
	case http.MethodGet:
		handleSearchWord(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func handleStoreWord(w http.ResponseWriter, r *http.Request) {
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
	prefix := r.FormValue("prefix")
	if !isValidWord(prefix) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Invalid format")
		return
	}

	prefix = strings.ToLower(prefix)
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

	if maxCount == 0 {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, "No matching word found in storage")
		return
	}
}

