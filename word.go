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