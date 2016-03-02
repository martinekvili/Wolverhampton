package main

import (
	"fmt"
	"log"
	"net/http"
)

// SSEventSource is an event source for a server job.
type SSEventSource struct {

	// A queue in which we store the messages that have already arrived,
	// so the newly subscribed clients never miss anything that happened in the past.
	messageQueue []string

	// Create a map of clients, the keys of the map are the channels
	// over which we can push messages to attached clients.  (The values
	// are just booleans and are meaningless. - This is a trick to use the go map as a set.)
	//
	clients map[chan string]bool

	// Channel into which new clients can be pushed
	//
	newClients chan chan string

	// Channel into which disconnected clients should be pushed
	//
	defunctClients chan chan string

	// A chanel on which we can signal the EventSource to close itself.
	closeSignal chan bool

	// Channel into which messages are pushed to be broadcast out
	// to attahed clients.
	//
	messages chan string
}

// CreateSSEventSource is the constructor for the SSEventSource.
func CreateSSEventSource() *SSEventSource {
	return &SSEventSource{
		messageQueue:   make([]string, 0),
		clients:        make(map[chan string]bool),
		newClients:     make(chan (chan string)),
		defunctClients: make(chan (chan string)),
		closeSignal:    make(chan bool),
		messages:       make(chan string),
	}
}

// Start handles the addition & removal of clients, as well as the broadcasting
// of messages out to clients that are currently attached.
//
func (b *SSEventSource) Start() {

	// Start a goroutine
	//
	go func() {

		// Loop endlessly
		//
		for {

			// Block until we receive from one of the
			// three following channels.
			select {

			case s := <-b.newClients:

				// There is a new client attached and we
				// want to start sending them messages.
				b.clients[s] = true
				log.Println("Added new client")

				for _, msg := range b.messageQueue {
					s <- msg
				}

			case s := <-b.defunctClients:

				// A client has dettached and we want to
				// stop sending them messages.
				delete(b.clients, s)
				close(s)

				log.Println("Removed client")

			case msg := <-b.messages:

				b.messageQueue = append(b.messageQueue, msg)

				// There is a new message to send.  For each
				// attached client, push the new message
				// into the client's message channel.
				for s := range b.clients {
					s <- msg
				}
				log.Printf("Broadcast message to %d clients", len(b.clients))

			case <-b.closeSignal:

				for s := range b.clients {
					delete(b.clients, s)
					close(s)
				}

				log.Println("Closed event source.")
				break
			}
		}
	}()
}

func (b *SSEventSource) AddClient() chan string {
	// Create a new channel, over which the event source can
	// send this client messages.
	messageChan := make(chan string)

	// Add this client to the map of those that should
	// receive updates
	b.newClients <- messageChan

	return messageChan
}

// This SSEventSource method handles and HTTP requests.
func (b *SSEventSource) ServeHTTP(w http.ResponseWriter, messageChan chan string) {
	// Make sure that the writer supports flushing.
	f, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	// Listen to the closing of the http connection via the CloseNotifier
	notify := w.(http.CloseNotifier).CloseNotify()
	go func() {
		<-notify
		// Remove this client from the map of attached clients
		// when `EventHandler` exits.
		b.defunctClients <- messageChan
		log.Println("HTTP connection just closed.")
	}()

	// Set the headers related to event streaming.
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	for {

		// Read from our messageChan.
		msg, open := <-messageChan

		if !open {
			// If our messageChan was closed, this means that the client has
			// disconnected.
			break
		}

		// Write to the ResponseWriter, `w`.
		fmt.Fprintf(w, "data: Message: %s\n\n", msg)

		// Flush the response.  This is only possible if
		// the repsonse supports streaming.
		f.Flush()
	}

	// Done.
}
