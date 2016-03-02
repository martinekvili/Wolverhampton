package main

import (
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
	log.Printf("Added event source: %v\n", sourceNum)
}

func (b *SSEventBroker) GetEventSource(sourceNum int) *SSEventSource {
	b.syncObject.RLock()
	defer b.syncObject.RUnlock()

	log.Printf("Got event source: %v\n", sourceNum)
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

func (b *SSEventBroker) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	eventNum, err := strconv.Atoi(r.URL.Path[len("/events/"):])
	if err != nil {
		log.Printf("Couldn't parse %s\n", r.URL.Path[len("/events/"):])

		// Wrong URL parameter
		return
	}

	eventSource, messageChan := b.getMessageChannelForEventSource(eventNum)
	if eventSource != nil {
		eventSource.ServeHTTP(w, messageChan)
	}
}
