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

