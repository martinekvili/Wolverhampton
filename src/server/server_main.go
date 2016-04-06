package main

import (
	"log"
	"net"
	"net/http"
	"net/rpc"
)

func main() {
	serviceContract := new(ServiceContract)
	rpc.Register(serviceContract)
	rpc.HandleHTTP()

	l, e := net.Listen("tcp", ":1234")
	if e != nil {
		log.Printf("Error while trying to listen on tcp: %v\n", e)
	}
	http.Serve(l, nil)
}
