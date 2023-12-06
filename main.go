package main

import (
	"fmt"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/sse-test", eventsHandler)
	http.ListenAndServe(":8080", nil)
}

func eventsHandler(w http.ResponseWriter, r *http.Request) {
	// Set headers for Server-Sent Events
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-transform, no-cache")

	// Create a channel to send data to the client
	dataChannel := make(chan string)

	// Close the channel when the connection is closed
	defer close(dataChannel)

	// Send dummy data to the client in a second loop
	go func() {
		for {
			data := "data: " + time.Now().Format("15:04:05") + "\n\n"
			fmt.Println(data)
			dataChannel <- data
			time.Sleep(1 * time.Second)
		}
	}()

	// Continuously write data to the client
	for {
		select {
		case data := <-dataChannel:
			fmt.Fprint(w, data)
			w.(http.Flusher).Flush()
		case <-r.Context().Done():
			return
		}
	}
}
