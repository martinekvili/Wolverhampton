package main

import (
	"datacontract"
	"fmt"
)

type CallbackContract bool

func (c *CallbackContract) SendBuildResult(args *datacontract.BuildResultArgs, resp *datacontract.EmptyArgs) error {
	eventBroker := GetSSEventBrokerInstance()
	eventBroker.GetEventSource(args.JobID).messages <- fmt.Sprintf("The build succeeded: %v", args.BuildResult)
	return nil
}

func (c *CallbackContract) CloseJob(args *datacontract.CloseJobArgs, resp *datacontract.EmptyArgs) error {
	eventBroker := GetSSEventBrokerInstance()
	eventBroker.RemoveEventSource(args.JobID)
	return nil
}

func (c *CallbackContract) JobStatus(args *datacontract.JobStatusArgs, resp *datacontract.EmptyArgs) error {
	eventBroker := GetSSEventBrokerInstance()
	eventBroker.GetEventSource(args.JobID).messages <- fmt.Sprintf("The job is %v. in line.", args.JobNumInRow+1)
	return nil
}

func (c *CallbackContract) NextJobStarted(args *datacontract.EmptyArgs, resp *datacontract.EmptyArgs) error {
	eventBroker := GetSSEventBrokerInstance()
	eventBroker.BroadCastEvent("The next job has started.")
	return nil
}
