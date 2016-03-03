package main

import (
	"datacontract"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/rpc"
	"strconv"
	"time"
)

// Page is the model for the html page
type Page struct {
	Name           string
	Color          string
	EventSourceNum int
}

func handler(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	client, err := rpc.DialHTTP("tcp", "localhost:1234")
	if err != nil {
		fmt.Printf("dialing: %v\n", err)
	}

	var args datacontract.EmptyArgs
	args.Name = "Vili"
	args.Number = 1234
	var reply datacontract.EmptyResult
	err = client.Call("Empty.DoSomething", args, &reply)
	if err != nil {
		fmt.Printf("arith error: %v\n", err)
	}

	fmt.Fprintf(w, "When started, the time was %v,\nat the end it is %v.\nResult is: %s\n", startTime, time.Now(), reply.Result)
}

// MainPageHandler is the handler for the main page, which we wire up to the
// route at "/" below in `main`.
//
func MainPageHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Got in http request")

	sourceNum, err := strconv.Atoi(r.URL.Path[1:])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	eventBroker := GetSSEventBrokerInstance()

	if eventBroker.HasEventSource(sourceNum) {
		// Read in the template with our SSE JavaScript code.
		t, err := template.ParseFiles("../src/webserver/templates/index.html")
		if err != nil {
			log.Fatal("WTF dude, error parsing your template.")

		}

		// Render the template, writing to `w`.
		var p Page
		p.Name = "Duder"
		p.Color = "Green"
		p.EventSourceNum = sourceNum
		t.Execute(w, p)
	} else {
		fmt.Fprintf(w, "<html><body>It'll be read from database</body></html>")
	}

	// Done.
	log.Println("Finished HTTP request at ", r.URL.Path)
}

// Main routine
//
func main() {

	// Make a new Broker instance
	eventBroker := GetSSEventBrokerInstance()

	// Start processing events
	eventBroker.AddEventSource(1234)

	// Make b the HTTP handler for "/events/".  It can do
	// this because it has a ServeHTTP method.  That method
	// is called in a separate goroutine for each
	// request to "/events/".
	http.Handle("/events/", eventBroker)

	// Generate a constant stream of events that get pushed
	// into the Broker's messages channel and are then broadcast
	// out to any clients that are attached.
	go func() {
		for i := 0; ; i++ {

			// Create a little message to send to clients,
			// including the current time.
			eventBroker.GetEventSource(1234).messages <- fmt.Sprintf("%d - the time is %v", i, time.Now())

			// Print a nice log message and sleep for 5s.
			log.Printf("Sent message %d ", i)
			time.Sleep(5 * 1e9)

			if i >= 5 {
				break
			}
		}

		eventBroker.RemoveEventSource(1234)
	}()

	// When we get a request at "/", call `MainPageHandler`
	// in a new goroutine.
	http.Handle("/", http.HandlerFunc(MainPageHandler))

	// Start the server and listen forever on port 8000.
	http.ListenAndServe(":8000", nil)
}
