package main

import (
	"datacontract"
	"fmt"
	"net"
	"net/http"
	"net/rpc"
	"sync/atomic"
	"time"
)

var cntr int32

type Empty bool

func (e *Empty) DoSomething(args *datacontract.EmptyArgs, resp *datacontract.EmptyResult) error {
	tmpCntr := atomic.AddInt32(&cntr, 1)

	fmt.Printf("Entered rpc - No. %v\n", tmpCntr)
	dur, _ := time.ParseDuration("3s")
	time.Sleep(dur)

	resp.Result = fmt.Sprintf("Hello there %v#%v, it's always a pleasure to see you!", args.Name, args.Number)

	fmt.Printf("Exited rpc - No. %v\n", tmpCntr)
	return nil
}

func main() {
	empty := new(Empty)
	rpc.Register(empty)
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", ":1234")
	if e != nil {
		fmt.Printf("listen error: %v\n", e)
	}
	http.Serve(l, nil)
}
