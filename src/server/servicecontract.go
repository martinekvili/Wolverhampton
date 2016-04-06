package main

import (
    "datacontract"
    "sync/atomic"
    "encoding/xml"
	"io/ioutil"
	"os/exec"
    "log"
    "errors"
    "fmt"
    "time"
)

var cntr int32 = 1

type ServiceContract bool

func (s *ServiceContract) GetID(args *datacontract.EmptyArgs, resp *int) error {
	*resp = int(atomic.AddInt32(&cntr, 1))
	return nil
}

func (s *ServiceContract) BuildProject(args *string, resp *bool) error {
    cmd := exec.Command("C:\\Program Files (x86)\\MSBuild\\12.0\\Bin\\amd64\\MSBuild.exe",
		"/noconsolelogger",
		"/logger:E:\\GitHub\\Wolverhampton\\XmlLogger\\XmlLogger\\bin\\Debug\\XmlLogger.dll",
		"TestSolution.sln")
	cmd.Dir = "E:\\BME\\onlab\\TestSolution\\"

	cmd.Run() // Although this can return us an error, it is not really usable
	// since this gives an error when the compilation process succeeded
	// but the compiled code had errors.
	// That's why we only check for the existence of the buildresult.xml file.

	bytes, err := ioutil.ReadFile("E:\\BME\\onlab\\TestSolution\\buildresult.xml")

	if err != nil {
        errorMsg := "There was an error with the compilation."
		log.Println(errorMsg)
		return errors.New(errorMsg)
	}

	var v BuildResult
	err2 := xml.Unmarshal(bytes, &v)
	if err2 != nil {
        errorMsg := fmt.Sprintf("Error happened during unmarshalling xml: %v", err2)
		log.Println(errorMsg)
		return errors.New(errorMsg)
	}
    
    time.Sleep(5 * time.Second)

	*resp = v.Successful
    return nil
}