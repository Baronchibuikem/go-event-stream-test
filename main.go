package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

const publicAPIURL = "https://jsonplaceholder.typicode.com/posts"

func streamAPIResponses(w http.ResponseWriter, r *http.Request) {
	// Set headers for SSE
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	// Make a request to the public API
	resp, err := http.Get(publicAPIURL)
	if err != nil {
		log.Println("Error fetching data from API:", err)
		http.Error(w, "Unable to fetch data", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Decode the JSON response
	var posts []map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&posts)
	if err != nil {
		log.Println("Error decoding API response:", err)
		http.Error(w, "Invalid API response", http.StatusInternalServerError)
		return
	}

	// Stream the posts one by one
	for i, post := range posts {
		// Format the post as a JSON string
		postData, err := json.Marshal(post)
		if err != nil {
			log.Println("Error marshaling post:", err)
			break
		}

		// Write the event to the response
		event := fmt.Sprintf("data: %s\n\n", postData)
		_, err = fmt.Fprint(w, event)
		if err != nil {
			log.Println("Error writing to client:", err)
			break
		}

		// Flush the response
		flusher, ok := w.(http.Flusher)
		if ok {
			flusher.Flush()
		}

		// Simulate delay between events
		time.Sleep(1 * time.Second)

		// Stop after 10 posts for demonstration purposes
		if i >= 9 {
			break
		}
	}
}

func main() {
	http.HandleFunc("/events", streamAPIResponses)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	log.Println("Server running on http://localhost:3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}
