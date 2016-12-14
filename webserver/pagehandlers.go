package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/rpc"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strconv"

	"github.com/martinekvili/Wolverhampton/datacontract"
)

func checkSession(w http.ResponseWriter, r *http.Request) bool {
	cookie, err := r.Cookie("session")

	if err == nil && cookie.Value != "" {
		sessionID, err := url.QueryUnescape(cookie.Value)

		if err == nil {
			hasSession, _ := GetSessionHandlerInstance().GetUserType(sessionID)

			if hasSession {
				return true
			}
		}
	}

	t, err := template.ParseFiles("templates/login.html")
	if err != nil {
		log.Fatal("Bad template for login.html")
		return false
	}

	t.Execute(w, nil)
	return false
}

// IndexPageHandler handles the index.html page
func IndexPageHandler(w http.ResponseWriter, r *http.Request) {
	if !checkSession(w, r) {
		return
	}

	t, err := template.ParseFiles("templates/index.html")
	if err != nil {
		log.Fatal("Bad template for index.html")
		return
	}

	t.Execute(w, nil)
}

// StaticFilesHandler handles the static file requests
func StaticFilesHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, r.URL.Path[1:])
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var credentials datacontract.LoginCredentials
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		log.Print(err)
	}
	if err := r.Body.Close(); err != nil {
		log.Print(err)
	}
	if err := json.Unmarshal(body, &credentials); err != nil {
		log.Print(err)
	}

	client, err := rpc.DialHTTP("tcp", "localhost:1234")
	if err != nil {
		log.Printf("Error happened while dialing: %v\n", err)
		return
	}

	var result datacontract.LoginResponse
	err = client.Call("ServiceContract.LoginUser", credentials, &result)
	if err != nil {
		log.Printf("Error happened during remote procedure call: %v\n", err)
		return
	}

	if !result.Success {
		json.NewEncoder(w).Encode(false)
	} else {
		sessionID := GetSessionHandlerInstance().CreateSession(result.User.Name, result.User.UserType)
		cookie := http.Cookie{Name: "session", Value: url.QueryEscape(sessionID), Path: "/", HttpOnly: true, MaxAge: 365 * 24 * 60 * 60}
		http.SetCookie(w, &cookie)
		json.NewEncoder(w).Encode(true)
	}
}

func ListUsersHandler(w http.ResponseWriter, r *http.Request) {
	client, err := rpc.DialHTTP("tcp", "localhost:1234")
	if err != nil {
		log.Printf("Error happened while dialing: %v\n", err)
		return
	}

	var args datacontract.EmptyArgs
	var result datacontract.UserList
	err = client.Call("ServiceContract.ListUsers", args, &result)
	if err != nil {
		log.Printf("Error happened during remote procedure call: %v\n", err)
		return
	}

	var users []UserViewModel
	for _, u := range result.Users {
		var userType string
		if u.UserType == datacontract.Admin {
			userType = "Admin"
		} else if u.UserType == datacontract.Teacher {
			userType = "Teacher"
		} else {
			userType = "Student"
		}

		users = append(users, UserViewModel{
			UserName: u.Name,
			FullName: u.FullName,
			UserType: userType,
		})
	}

	json.NewEncoder(w).Encode(users)
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
	var jobID int
	err = client.Call("ServiceContract.GetID", args, &jobID)
	if err != nil {
		log.Printf("Error happened during remote procedure call: %v\n", err)
		return
	}

	eventBroker := GetSSEventBrokerInstance()
	eventBroker.AddEventSource(jobID)

	http.Redirect(w, r, "/id/"+strconv.Itoa(jobID), 302)
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

		fileName := path.Join("uploads", strconv.Itoa(jobID)+filepath.Ext(handler.Filename))
		os.MkdirAll(path.Dir(fileName), os.ModeDir)

		f, err := os.Create(fileName)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer f.Close()

		io.Copy(f, file)

		startJobArgs := &datacontract.StartJobArgs{
			JobID:    jobID,
			FileName: handler.Filename,
		}
		var reply datacontract.EmptyArgs
		client.Call("ServiceContract.StartJob", startJobArgs, &reply)
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
		t, err := template.ParseFiles("templates/testpage.html")
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
		client, err := rpc.DialHTTP("tcp", "localhost:1234")
		if err != nil {
			log.Printf("Error happened while dialing: %v\n", err)
			return
		}

		getJobResultArgs := &datacontract.GetJobResultArgs{
			JobID: sourceNum,
		}
		var jobResult datacontract.JobResult
		err = client.Call("ServiceContract.GetJobResult", getJobResultArgs, &jobResult)
		if err != nil {
			log.Printf("Error happened during remote procedure call: %v\n", err)
			return
		}

		fmt.Fprintf(w, "<html><body>The build was %v ly succesful.</body></html>", jobResult.BuildInfo.Successful)
	}

	// Done.
	log.Println("Finished HTTP request at ", r.URL.Path)
}
