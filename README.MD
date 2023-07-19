# instructions

Write an http service that will accept 2 endpoints:

1. The first endpoint will accept a word (eg. "animal", "house") and store it in the service. Definition for a word is the following regex `[a-zA-Z]+`
2. The second endpoint will accept the beginning of a word (eg. "an") and returns the most frequent word stored in the service
3. The service will be case insensitive

scenario example
```
POST /service/word="abc"
POST /service/word="ab"
POST /service/word="ab"

GET /service/prefix="a"     => response: "ab"
GET /service/prefix="ab"   => response: "ab"
GET /service/prefix="abc"   => response: "abc"
GET /service/prefix="d"     => response: null
```

**Its up to you to choose a design in the API that makes sense**

The service must be written in Go, expect that there is no restrictions in terms of technologies or designs

Bonus points:
- Use docker
- Scalability
- Performances

curl -X POST "http://localhost:8080/store" -d "word=animal"
curl "http://localhost:8080/search?prefix=an"


This program is a Go application that implements a simple server with two endpoints using the Gin web framework. The purpose of this program is to allow users to store words and then search for the most frequently stored word with a specific prefix.

Defines a global variable called "storage," which is the map used to store words and their counts. Also, it utilises a "storageLock" variable, which is a sync.RWMutex to handle concurrent access to the "storage" map.

The "setupRouter" function sets up the Gin router with two endpoints: "/store" and "/search." The "/store" endpoint is used to store words, and the "/search" endpoint is used to search for the most frequently stored word with a specific prefix.

The "/store" endpoint handles POST requests and expects a word to be sent in the request body. It checks if the word is in the correct format (containing only letters), converts it to lowercase (to make this case insensitive as per requirements), and then increments its count in the "storage" map. After successful storage, it responds with a JSON message indicating the word has been stored.

The "/search" endpoint handles GET requests and expects a "prefix" query parameter to be provided. It first validates the format of the prefix (letters only), converts it to lowercase, and then calls the "getMaxWordWithPrefixConcurrent" function to find the most frequently stored word with the given prefix.

The "getMaxWordWithPrefixConcurrent" function performs a concurrent search through the "storage" map for words with the specified prefix. It utilises goroutines to parallelise the search process (meaning multiple goroutines are executing simultaneously). Once all goroutines have completed, it identifies the word with the highest frequency and returns it as the result.

The "isValidWord" function uses a regular expression to check if a given word contains only letters (no digits, spaces, or special characters).

In the "main" function, the Gin router is initialised, and the server starts listening on "localhost:8080" using the "http.ListenAndServe" function. The server runs in a goroutine to enable concurrent handling of requests. The program enters an infinite loop using "select {}" to keep the server running indefinitely.

Basically, this programme provides a basic HTTP server that allows users to store words and find the most frequently stored word with a specified prefix.
