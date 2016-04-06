package main

import (
	"net/http"
)

// Main routine
//
func main() {
	// Make a new Broker instance
	eventBroker := GetSSEventBrokerInstance()

	http.Handle("/events/", eventBroker)
	http.Handle("/id/", http.HandlerFunc(TestPageHandler))
	http.Handle("/submit/", http.HandlerFunc(SubmitPageHandler))
	http.Handle("/", http.HandlerFunc(IndexPageHandler))

	// Start the server and listen forever on port 8000.
	http.ListenAndServe(":8000", nil)
}
