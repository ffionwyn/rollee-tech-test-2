package main

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
)

var (
	storage     = make(map[string]int)
	storageLock sync.RWMutex
)

func main() {
	router := gin.Default()

	router.POST("/store", handleStoreWord)
	router.GET("/search", handleSearchWord)

	log.Fatal(router.Run(":8080"))
}

func handleStoreWord(c *gin.Context) {
	word := c.PostForm("word")
	if !isValidWord(word) {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "incorrect format"})
		return
	}

	word = strings.ToLower(word)
	storageLock.Lock()
	defer storageLock.Unlock()
	storage[word]++

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("'%s' stored successfully", word)})
}

func handleSearchWord(c *gin.Context) {
	prefix := c.DefaultQuery("prefix", "")
	if !isValidWord(prefix) {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid prefix format"})
		return
	}

	prefix = strings.ToLower(prefix)
	maxWord := getMaxWordWithPrefixConcurrent(prefix)

	if maxWord == "" {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "No matching word found in the storage"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Most frequent word with prefix '%s': %s", prefix, maxWord)})
	fmt.Println("Result:", maxWord)
}

func getMaxWordWithPrefixConcurrent(prefix string) string {
	type result struct {
		word  string
		count int
	}

	results := make(chan result, len(storage))
	var wg sync.WaitGroup

	storageLock.RLock()
	for word, count := range storage {
		wg.Add(1)
		go func(word string, count int) {
			defer wg.Done()
			wordLower := strings.ToLower(word)
			if strings.HasPrefix(wordLower, prefix) {
				results <- result{word: word, count: count}
			}
		}(word, count)
	}
	storageLock.RUnlock()

	wg.Wait()
	close(results)

	var maxWord string
	maxCount := 0
	for res := range results {
		if res.count > maxCount {
			maxWord = res.word
			maxCount = res.count
		}
	}

	return maxWord
}

func isValidWord(word string) bool {
	validWordRegex := regexp.MustCompile(`^([a-zA-Z])+$`)
	return validWordRegex.MatchString(word)
}