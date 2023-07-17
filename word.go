package main

import (
	"log"
	"net/http"
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

