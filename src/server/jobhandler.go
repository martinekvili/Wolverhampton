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
	"strings"
	"time"
)

func createDirectoryForJob(jobID int) (jobStoragePath string, solutionPath string) {
	jobStoragePath = path.Join("JobStorage", strconv.Itoa(jobID))
	os.MkdirAll(jobStoragePath, os.ModeDir)

	ExtractZipIntoFolder(path.Join("projects", "TestSolution.zip"), jobStoragePath)
	solutionPath = path.Join(jobStoragePath, "TestSolution")
	CopyFile(path.Join("uploads", strconv.Itoa(jobID)+".cs"), path.Join(solutionPath, "TestSolution", "ClassToBeWritten.cs"))

	return
}

func buildProject(jobID int, solutionPath string, client *rpc.Client) error {
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
		return err
	}

	var v BuildResult
	err = xml.Unmarshal(bytes, &v)
	if err != nil {
		errorMsg := fmt.Sprintf("Error happened during unmarshalling xml: %v", err)
		log.Println(errorMsg)
		return err
	}

	buildArgs := &datacontract.BuildResultArgs{
		JobID:       jobID,
		BuildResult: v.Successful,
	}
	var reply datacontract.EmptyArgs
	client.Call("CallbackContract.SendBuildResult", buildArgs, &reply)

	return nil
}

func runProject(jobID int, solutionPath string, client *rpc.Client) error {
	cmd := exec.Command("E:\\GitHub\\Wolverhampton\\SandBoxRunner\\Debug\\SandboxRunner.exe",
		"5",
		"5",
		"TestSolution\\bin\\Debug\\TestSolution.exe",
		"runresult.txt")
	cmd.Dir = solutionPath

	outputBytes, err := cmd.Output()
	if err != nil {
		log.Printf("Error happened during running the executable: %v\n", err)
		return err
	}

	output := string(outputBytes)
	log.Printf("Run result is: %v\n", output)

	var runResult datacontract.RunResult
	switch {
	case strings.Contains(output, "SUCCESS"):
		runResult = datacontract.Success

	case strings.Contains(output, "NOT_ENOUGH_MEMORY"):
		runResult = datacontract.NotEnoughMemory

	case strings.Contains(output, "NOT_ENOUGH_TIME"):
		runResult = datacontract.NotEnoughTime

	case strings.Contains(output, "UNKNOWN_ERROR"):
		runResult = datacontract.Unknown
	}

	runResultArgs := &datacontract.RunResultArgs{
		JobID:  jobID,
		Result: runResult,
	}
	var reply datacontract.EmptyArgs
	client.Call("CallbackContract.SendRunResult", runResultArgs, &reply)

	return nil
}

func closeJob(jobID int, client *rpc.Client) {
	closeArgs := &datacontract.CloseJobArgs{
		JobID: jobID,
	}
	var reply datacontract.EmptyArgs
	client.Call("CallbackContract.CloseJob", closeArgs, &reply)
}

func HandleJob(jobID int, doneWork chan bool) {
	/*jobStoragePath*/ _, solutionPath := createDirectoryForJob(jobID)
	//defer os.RemoveAll(jobStoragePath)

	client, err := rpc.DialHTTP("tcp", "localhost:1235")
	if err != nil {
		log.Printf("Error happened while dialing: %v\n", err)
		return
	}

	err = buildProject(jobID, solutionPath, client)
	if err != nil {
		return
	}

	err = runProject(jobID, solutionPath, client)
	if err != nil {
		return
	}

	time.Sleep(time.Second * 10)

	closeJob(jobID, client)

	doneWork <- true
}
