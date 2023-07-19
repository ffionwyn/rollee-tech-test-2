package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)


func TestHandleStoreWord(t *testing.T) {
	router := gin.Default() // new gin router

	req, _ := http.NewRequest("POST", "/store", nil)
	req.Form = map[string][]string{
		"word": {"dog"},
	}

	w := httptest.NewRecorder() // perform the request and record the response
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, w.Code)
	}

	expectedMessage := "'dog' stored successfully"
	if body := w.Body.String(); body != expectedMessage {
		t.Errorf("expected response body '%s', but got '%s'", expectedMessage, body)
	}
}

func TestHandleSearchWord(t *testing.T) {
	router := gin.Default() // new gin router

	storage = map[string]int{ // fake storing some words
		"dog":  3,
		"house": 5,
	}

	req, _ := http.NewRequest("GET", "/search?prefix=do", nil) // create a test request with "prefix" query parameter

	w := httptest.NewRecorder() // perform the request and record the response
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status code %d, but got %d", http.StatusOK, w.Code)
	}

	expectedMessage := "most frequent word with prefix 'do': dog"
	if body := w.Body.String(); body != expectedMessage {
		t.Errorf("expected response body '%s', but got '%s'", expectedMessage, body)
	}
}

func TestIsValidWord(t *testing.T) {
	if !isValidWord("hello") { // testing valid words
		t.Error("Expected 'hello' to be a valid word, but got invalid")
	}
	if !isValidWord("SHORTS") {
		t.Error("Expected 'SHORTS' to be a valid word, but got invalid")
	}

	if isValidWord("123") { // testing invalid words
		t.Error("Expected '123' to be an invalid word, but got valid")
	}
	if isValidWord("hello-world") {
		t.Error("Expected 'hello-world' to be an invalid word, but got valid")
	}
}


