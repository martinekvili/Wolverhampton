package main

import (
	"datacontract"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/rpc"
	"os"
	"os/exec"
	"path"
	"strconv"
	"time"
)

func HandleJob(jobID int, doneWork chan bool) {
	jobStoragePath := path.Join("JobStorage", strconv.Itoa(jobID))
	os.MkdirAll(jobStoragePath, os.ModeDir)

	ExtractZipIntoFolder(path.Join("projects", "TestSolution.zip"), jobStoragePath)
	solutionPath := path.Join(jobStoragePath, "TestSolution")
	CopyFile(path.Join("uploads", strconv.Itoa(jobID)+".cs"), path.Join(solutionPath, "TestSolution", "ClassToBeWritten.cs"))

	cmd := exec.Command("C:\\Program Files (x86)\\MSBuild\\12.0\\Bin\\amd64\\MSBuild.exe",
		"/noconsolelogger",
		"/logger:E:\\GitHub\\Wolverhampton\\XmlLogger\\XmlLogger\\bin\\Debug\\XmlLogger.dll",
		"TestSolution.sln")
	cmd.Dir = solutionPath

	cmd.Run() // Although this can return us an error, it is not really usable
	// since this gives an error when the compilation process succeeded
	// but the compiled code had errors.
	// That's why we only check for the existence of the buildresult.xml file.

	bytes, err := ioutil.ReadFile(path.Join(solutionPath, "buildresult.xml"))

	if err != nil {
		errorMsg := "There was an error with the compilation."
		log.Println(errorMsg)
		return
	}

	var v BuildResult
	err = xml.Unmarshal(bytes, &v)
	if err != nil {
		errorMsg := fmt.Sprintf("Error happened during unmarshalling xml: %v", err)
		log.Println(errorMsg)
		return
	}

	client, err := rpc.DialHTTP("tcp", "localhost:1235")
	if err != nil {
		log.Printf("Error happened while dialing: %v\n", err)
		return
	}

	buildArgs := &datacontract.BuildResultArgs{
		JobID:       jobID,
		BuildResult: v.Successful,
	}
	var reply datacontract.EmptyArgs
	err = client.Call("CallbackContract.SendBuildResult", buildArgs, &reply)

	time.Sleep(time.Second * 10)

	closeArgs := &datacontract.CloseJobArgs{
		JobID: jobID,
	}
	err = client.Call("CallbackContract.CloseJob", closeArgs, &reply)

	doneWork <- true
}
