package main

import (
	"bufio"
	"datacontract"
	"encoding/xml"
	"fmt"
	"gopkg.in/mgo.v2"
	"io/ioutil"
	"log"
	"net/rpc"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
)

func createDirectoryForJob(jobID int) (jobStoragePath string, solutionPath string) {
	jobStoragePath = path.Join("JobStorage", strconv.Itoa(jobID))
	os.MkdirAll(jobStoragePath, os.ModeDir)

	ExtractZipIntoFolder(path.Join("projects", "TestSolution.zip"), jobStoragePath)
	solutionPath = path.Join(jobStoragePath, "TestSolution")
	CopyFile(path.Join("uploads", strconv.Itoa(jobID)+".cs"), path.Join(solutionPath, "TestSolution", "ClassToBeWritten.cs"))

	return
}

func buildProject(jobID int, solutionPath string, jobResult *datacontract.JobResult, client *rpc.Client) error {
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

	var result datacontract.BuildResult
	err = xml.Unmarshal(bytes, &result)
	if err != nil {
		errorMsg := fmt.Sprintf("Error happened during unmarshalling xml: %v", err)
		log.Println(errorMsg)
		return err
	}

	jobResult.BuildInfo = result

	buildArgs := &datacontract.BuildResultArgs{
		JobID:  jobID,
		Result: result,
	}
	var reply datacontract.EmptyArgs
	client.Call("CallbackContract.SendBuildResult", buildArgs, &reply)

	return nil
}

func runProject(jobID int, solutionPath string, jobResult *datacontract.JobResult, client *rpc.Client) (datacontract.RunResult, error) {
	cmd := exec.Command("E:\\GitHub\\Wolverhampton\\SandBoxRunner\\Debug\\SandboxRunner.exe",
		"10",
		"5",
		"TestSolution\\bin\\Debug\\TestSolution.exe",
		"output.txt")
	cmd.Dir = solutionPath

	outputBytes, err := cmd.Output()
	if err != nil {
		log.Printf("Error happened during running the executable: %v\n", err)
		return datacontract.Unknown, err
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

	jobResult.RunInfo = runResult

	runResultArgs := &datacontract.RunResultArgs{
		JobID:  jobID,
		Result: runResult,
	}
	var reply datacontract.EmptyArgs
	client.Call("CallbackContract.SendRunResult", runResultArgs, &reply)

	return runResult, nil
}

func matchOutput(jobID int, solutionPath string, jobResult *datacontract.JobResult, client *rpc.Client) error {
	expectedOutput, err := os.Open(path.Join(solutionPath, "output_expected.txt"))
	if err != nil {
		log.Printf("Error happened during opening expected output file: %v\n", err)
		return err
	}
	defer expectedOutput.Close()

	output, err := os.Open(path.Join(solutionPath, "output.txt"))
	if err != nil {
		log.Printf("Error happened during opening output file: %v\n", err)
		return err
	}
	defer output.Close()

	resultArgs := &datacontract.OutputMatchResultArgs{
		JobID:      jobID,
		Mismatches: make([]datacontract.OutputMismatchLine, 0),
	}

	expectedOutputScanner := bufio.NewScanner(expectedOutput)
	outputScanner := bufio.NewScanner(output)

	lineNum := 0
	for expectedOutputScanner.Scan() && outputScanner.Scan() {
		expectedLine := expectedOutputScanner.Text()
		line := outputScanner.Text()

		if expectedLine != line {
			mismatch := &datacontract.OutputMismatchLine{
				LineNumber: lineNum,
				Expected:   expectedLine,
				Actual:     line,
			}

			resultArgs.Mismatches = append(resultArgs.Mismatches, *mismatch)
		}

		lineNum++
	}

	for expectedOutputScanner.Scan() {
		expectedLine := expectedOutputScanner.Text()

		mismatch := &datacontract.OutputMismatchLine{
			LineNumber: lineNum,
			Expected:   expectedLine,
			Actual:     "",
		}

		resultArgs.Mismatches = append(resultArgs.Mismatches, *mismatch)

		lineNum++
	}

	for outputScanner.Scan() {
		line := outputScanner.Text()

		mismatch := &datacontract.OutputMismatchLine{
			LineNumber: lineNum,
			Expected:   "",
			Actual:     line,
		}

		resultArgs.Mismatches = append(resultArgs.Mismatches, *mismatch)

		lineNum++
	}

	jobResult.CompareInfo = resultArgs.Mismatches

	var reply datacontract.EmptyArgs
	client.Call("CallbackContract.SendOutputMatchResult", resultArgs, &reply)

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
	defer func() { doneWork <- true }()

	jobResult := &datacontract.JobResult{
		JobID: jobID,
		BuildInfo: datacontract.BuildResult{
			Successful: false,
			ErrorList:  make([]datacontract.Error, 0),
		},
		RunInfo:     datacontract.Unknown,
		CompareInfo: make([]datacontract.OutputMismatchLine, 0),
	}

	defer func() {
		session, err := mgo.Dial("localhost")
		if err != nil {
			log.Printf("Couldn't access database: %v\n", err)
		}

		c := session.DB("WolverhamptonDB").C("JobResult")
		err = c.Insert(jobResult)
		if err != nil {
			log.Printf("Couldn't insert job result: %v\n", err)
		}

		session.Close()
	}()

	jobStoragePath, solutionPath := createDirectoryForJob(jobID)
	defer os.RemoveAll(jobStoragePath)

	client, err := rpc.DialHTTP("tcp", "localhost:1235")
	if err != nil {
		log.Printf("Error happened while dialing: %v\n", err)
		return
	}
	defer closeJob(jobID, client)

	err = buildProject(jobID, solutionPath, jobResult, client)
	if err != nil {
		return
	}

	runResult, err := runProject(jobID, solutionPath, jobResult, client)
	if err != nil || runResult != datacontract.Success {
		return
	}

	matchOutput(jobID, solutionPath, jobResult, client)
}
