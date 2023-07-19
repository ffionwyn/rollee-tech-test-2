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

var maxWord string

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
	maxCount := 0

	storageLock.RLock()
	defer storageLock.RUnlock()
	for word, count := range storage {
		wordLower := strings.ToLower(word)
		if strings.HasPrefix(wordLower, prefix) && count > maxCount {
			maxCount = count
			maxWord = word
		}
	}

	if maxCount == 0 {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "No matching word found in the storage"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Most frequent word with prefix '%s': %s", prefix, getMaxWordWithPrefix(prefix))})
	fmt.Println("Result:", maxWord)
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