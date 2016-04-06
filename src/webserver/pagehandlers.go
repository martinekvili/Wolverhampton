package main

import (
    "datacontract"
    "net/http"
    "net/rpc"
    "html/template"
    "log"
    "strconv"
    "fmt"
    "os"
    "io"
)

// IndexPageHandler handles the index.html page
func IndexPageHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("../src/webserver/templates/index.html")
	if err != nil {
		log.Fatal("Bad template for index.html")
        return
	}

	t.Execute(w, nil)
}

// SubmitPageHandler handles the submits, and redirects the client to the test page
func SubmitPageHandler(w http.ResponseWriter, r *http.Request) {
    log.Println("Got in submit")
    client, err := rpc.DialHTTP("tcp", "localhost:1234")
	if err != nil {
		log.Printf("Error happened while dialing: %v\n", err)
        return
	}

	var args datacontract.EmptyArgs
	var reply int
	err = client.Call("ServiceContract.GetID", args, &reply)
	if err != nil {
		log.Printf("Error happened during remote procedure call: %v\n", err)
        return
	}
    
    eventBroker := GetSSEventBrokerInstance()
	eventBroker.AddEventSource(reply)
    
    http.Redirect(w, r, "/id/" + strconv.Itoa(reply), 302)
    log.Println("Redirected")
    
    // Parse up to 32 MB
    r.ParseMultipartForm(32 << 20)
    
    file, handler, err := r.FormFile("uploadedFile")
    if err != nil {
        log.Fatal("Error opening uploaded file.")
        return
    }
    
    go func() {
        defer file.Close()
        
        f, err := os.Create("uploads/" + handler.Filename)
        if err != nil {
            fmt.Println(err)
            return
        }
        defer f.Close()
        
        io.Copy(f, file)
        
        var buildResult bool
        err = client.Call("ServiceContract.BuildProject", handler.Filename, &buildResult)
        if err != nil {
            log.Printf("Error happened during remote procedure call: %v\n", err)
            return
        }
        
        log.Printf("The build succeeded: %v\n", buildResult)
        eventBroker.GetEventSource(reply).messages <- fmt.Sprintf("The build succeeded: %v", buildResult)
        eventBroker.RemoveEventSource(reply)
    }()
}

// Page is the model for the html page
type Page struct {
	Name           string
	Color          string
	EventSourceNum int
}

// TestPageHandler handles the testing page
func TestPageHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Got in http request")

	sourceNum, err := strconv.Atoi(r.URL.Path[len("/id/"):])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	eventBroker := GetSSEventBrokerInstance()

	if eventBroker.HasEventSource(sourceNum) {
		// Read in the template with our SSE JavaScript code.
		t, err := template.ParseFiles("../src/webserver/templates/testpage.html")
		if err != nil {
			log.Fatal("Bad template for testpage.html")
            return
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