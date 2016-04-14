package main

import (
	"datacontract"
	"log"
	"net/rpc"
	"sync"
)

var instance *JobQueue
var once sync.Once

type JobQueue struct {
	addItem  chan int
	jobQueue *Queue
}

func GetJobQueueInstance() *JobQueue {
	once.Do(func() {
		instance = &JobQueue{
			addItem:  make(chan int),
			jobQueue: CreateQueue(),
		}
	})

	return instance
}

func (jq *JobQueue) Start() {
	go func() {
		jq.jobQueue.Start()

		client, err := rpc.DialHTTP("tcp", "localhost:1235")
		if err != nil {
			log.Printf("Error happened while dialing: %v\n", err)
			return
		}

		for {
			select {
			case item := <-jq.addItem:
				args := &datacontract.JobStatusArgs{
					JobID:       item,
					JobNumInRow: <-jq.jobQueue.getSize,
				}
				var reply datacontract.EmptyArgs
				client.Call("CallbackContract.JobStatus", args, &reply)

				jq.jobQueue.addItem <- item

			case item := <-jq.jobQueue.getItem:
				var args, reply datacontract.EmptyArgs
				client.Call("CallbackContract.NextJobStarted", args, &reply)

				HandleJob(item)
			}
		}
	}()
}
