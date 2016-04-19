package main

import (
	"datacontract"
	"sync/atomic"
	// "encoding/xml"
	// "io/ioutil"
	// "os/exec"
	"log"
	// "errors"
	// "fmt"
	// "time"
)

var cntr int32 = 0

type ServiceContract bool

func (s *ServiceContract) GetID(args *datacontract.EmptyArgs, resp *int) error {
	*resp = int(atomic.AddInt32(&cntr, 1))
	return nil
}

func (s *ServiceContract) StartJob(args *datacontract.StartJobArgs, resp *datacontract.EmptyArgs) error {
	log.Println("Starting job...")
	jobQueue := GetJobQueueInstance()
	jobQueue.addItem <- args.JobID
	return nil
}
