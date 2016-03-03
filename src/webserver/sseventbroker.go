package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
)

var instance *SSEventBroker
var once sync.Once

// SSEventBroker handles the event sources
type SSEventBroker struct {
	syncObject sync.RWMutex

	eventSources map[int]*SSEventSource
}

func GetSSEventBrokerInstance() *SSEventBroker {
	once.Do(func() {
		instance = &SSEventBroker{
			eventSources: make(map[int]*SSEventSource),
		}
	})

	return instance
}

// HasEventSource determines whether the broker has the specified event source.
func (b *SSEventBroker) HasEventSource(sourceNum int) bool {
	b.syncObject.RLock()
	defer b.syncObject.RUnlock()

	return b.eventSources[sourceNum] != nil
}

func (b *SSEventBroker) AddEventSource(sourceNum int) {
	b.syncObject.Lock()
	defer b.syncObject.Unlock()

	eventSource := CreateSSEventSource()
	eventSource.Start()
	b.eventSources[sourceNum] = eventSource
}

func (b *SSEventBroker) GetEventSource(sourceNum int) *SSEventSource {
	b.syncObject.RLock()
	defer b.syncObject.RUnlock()

	return b.eventSources[sourceNum]
}

func (b *SSEventBroker) RemoveEventSource(sourceNum int) {
	b.syncObject.Lock()
	defer b.syncObject.Unlock()

	b.eventSources[sourceNum].closeSignal <- true
	delete(b.eventSources, sourceNum)
}

func (b *SSEventBroker) getMessageChannelForEventSource(eventNum int) (*SSEventSource, chan string) {
	b.syncObject.RLock()
	defer b.syncObject.RUnlock()

	eventSource := b.eventSources[eventNum]
	if eventSource == nil {
		// Since we created the html page this event source has been closed.
		// Now we tell the browser about it, so it reloads the page.
		return nil, nil
	}

	return eventSource, eventSource.AddClient()
}

//func sendErrorMessage(w http.ResponseWriter)

func (b *SSEventBroker) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Make sure that the writer supports flushing.
	f, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	// Set the headers related to event streaming.
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	// Listen to the closing of the http connection via the CloseNotifier
	notify := w.(http.CloseNotifier).CloseNotify()

	eventNum, err := strconv.Atoi(r.URL.Path[len("/events/"):])
	if err != nil {
		log.Printf("Couldn't parse %s\n", r.URL.Path[len("/events/"):])

		// Wrong URL parameter
		fmt.Fprint(w, "data: ERROR\n\n")
		f.Flush()
		return
	}

	eventSource, messageChan := b.getMessageChannelForEventSource(eventNum)
	if eventSource == nil {
		fmt.Fprint(w, "data: ERROR\n\n")
		f.Flush()
		return
	}

	go func() {
		<-notify
		// Remove this client from the map of attached clients
		// when `EventHandler` exits.
		eventSource.defunctClients <- messageChan
		log.Println("HTTP connection just closed.")
	}()

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
