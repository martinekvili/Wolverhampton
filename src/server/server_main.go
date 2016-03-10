package main

import (
	"datacontract"
	"fmt"
	// "net"
	// "net/http"
	// "net/rpc"
	"encoding/xml"
	"io/ioutil"
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
	// empty := new(Empty)
	// rpc.Register(empty)
	// rpc.HandleHTTP()
	// l, e := net.Listen("tcp", ":1234")
	// if e != nil {
	// 	fmt.Printf("listen error: %v\n", e)
	// }
	// http.Serve(l, nil)

	bytes, err := ioutil.ReadFile("E:\\BME\\onlab\\build.xml")

	if err != nil {
		fmt.Print("No such file")
		return
	}

	var v BuildResult

	err2 := xml.Unmarshal(bytes, &v)
	if err2 != nil {
		fmt.Printf("error: %v", err2)
		return
	}

	if v.Successful {
		fmt.Println("Build SUCCESSFUL.\n")
	} else {
		fmt.Println("Build FAILED.\n")
	}

	for _, e := range v.ErrorList {
		fmt.Printf("%v in %v(%v, %v): %v - %v\n\n", e.Type, e.FileName, e.LineNumber, e.ColumnNumber, e.Code, e.Message)
	}
}
