package main

import (
	"log"
	"net"
	"net/http"
	"net/rpc"
)

// Main routine
//
func main() {
	// Set up the RPC callback interface
	go func() {
		callbackContract := new(CallbackContract)
		rpc.Register(callbackContract)
		rpc.HandleHTTP()

		l, e := net.Listen("tcp", ":1235")
		if e != nil {
			log.Printf("Error while trying to listen on tcp: %v\n", e)
		}
		http.Serve(l, nil)
	}()

	// Make a new Broker instance
	eventBroker := GetSSEventBrokerInstance()

	http.Handle("/api/users", http.HandlerFunc(ListUsersHandler))
	http.Handle("/events/", eventBroker)
	http.Handle("/id/", http.HandlerFunc(TestPageHandler))
	http.Handle("/submit/", http.HandlerFunc(SubmitPageHandler))
	http.Handle("/stat/", http.HandlerFunc(StaticFilesHandler))
	http.Handle("/login", http.HandlerFunc(LoginHandler))
	http.Handle("/", http.HandlerFunc(IndexPageHandler))

	// Start the server and listen forever on port 8000.
	http.ListenAndServe(":8000", nil)
}
