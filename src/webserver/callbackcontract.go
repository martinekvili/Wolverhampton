package main

import (
	"datacontract"
	"fmt"
)

type CallbackContract bool

func (c *CallbackContract) SendOutputMatchResult(args *datacontract.OutputMatchResultArgs, resp *datacontract.EmptyArgs) error {
	eventBroker := GetSSEventBrokerInstance()

	eventBroker.GetEventSource(args.JobID).messages <- fmt.Sprintf("The output has %v wrong lines.", len(args.Mismatches))

	return nil
}

func (c *CallbackContract) SendRunResult(args *datacontract.RunResultArgs, resp *datacontract.EmptyArgs) error {
	eventBroker := GetSSEventBrokerInstance()

	var message string
	switch args.Result {
	case datacontract.Success:
		message = "The run succeeded."

	case datacontract.NotEnoughMemory:
		message = "The run was terminated because it tried to use more memory than permitted."

	case datacontract.NotEnoughTime:
		message = "The run was terminated because it tried to run for more time than permitted."

	case datacontract.Unknown:
		message = "The run was terminated due to unknown reasons."
	}

	eventBroker.GetEventSource(args.JobID).messages <- message
	return nil
}

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
	var msg string

	if args.JobNumInRow == 0 {
		msg = fmt.Sprint("The job has started.")
	} else {
		msg = fmt.Sprintf("The job is %v. in line.", args.JobNumInRow)
	}
	eventBroker := GetSSEventBrokerInstance()
	eventBroker.GetEventSource(args.JobID).messages <- msg
	return nil
}

func (c *CallbackContract) NextJobStarted(args *datacontract.EmptyArgs, resp *datacontract.EmptyArgs) error {
	eventBroker := GetSSEventBrokerInstance()
	eventBroker.BroadCastEvent("The next job has started.")
	return nil
}
