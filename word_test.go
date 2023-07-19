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


