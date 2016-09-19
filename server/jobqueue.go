package main

import (
	"log"
	"net/rpc"
	"sync"

	"github.com/martinekvili/Wolverhampton/datacontract"
)

func when(expression bool, channel chan int) chan int {
	if expression {
		return channel
	}

	return nil
}

var instance *JobQueue
var once sync.Once

type JobQueue struct {
	addItem  chan int
	jobQueue *Queue

	workInProgress bool
	doneWork       chan bool
}

func GetJobQueueInstance() *JobQueue {
	once.Do(func() {
		instance = &JobQueue{
			addItem:        make(chan int),
			jobQueue:       CreateQueue(),
			workInProgress: false,
			doneWork:       make(chan bool),
		}
	})

	return instance
}

func (jq *JobQueue) Start() {
	go func() {
		client, err := rpc.DialHTTP("tcp", "localhost:1235")
		if err != nil {
			log.Printf("Error happened while dialing: %v\n", err)
			return
		}

		jq.jobQueue.Start(client)

		for {
			select {
			case item := <-jq.addItem:
				args := &datacontract.JobStatusArgs{
					JobID:       item,
					JobNumInRow: <-jq.jobQueue.getSize + 1,
				}
				var reply datacontract.EmptyArgs
				client.Call("CallbackContract.JobStatus", args, &reply)

				jq.jobQueue.addItem <- item

			case <-jq.doneWork:
				jq.workInProgress = false

			case item := <-when(!jq.workInProgress, jq.jobQueue.getItem):
				args := &datacontract.JobStatusArgs{
					JobID:       item,
					JobNumInRow: 0,
				}
				var reply datacontract.EmptyArgs
				client.Call("CallbackContract.JobStatus", args, &reply)

				jq.workInProgress = true
				go HandleJob(item, jq.doneWork)
			}
		}
	}()
}
