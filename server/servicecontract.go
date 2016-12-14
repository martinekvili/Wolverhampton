package main

import (
	"log"
	"sync/atomic"

	"github.com/martinekvili/Wolverhampton/datacontract"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
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

func (s *ServiceContract) GetJobResult(args *datacontract.GetJobResultArgs, resp *datacontract.JobResult) error {
	session, err := mgo.Dial("localhost")
	if err != nil {
		log.Printf("Couldn't access database: %v\n", err)
		return err
	}
	defer session.Close()

	c := session.DB("WolverhamptonDB").C("JobResult")

	err = c.Find(bson.M{"jobid": args.JobID}).One(&resp)
	if err != nil {
		log.Printf("Couldn't read job result: %v\n", err)
		return err
	}

	return nil
}

func (s *ServiceContract) LoginUser(args *datacontract.LoginCredentials, resp *datacontract.LoginResponse) error {
	user, success := CheckPassword(args.UserName, args.Password)
	resp.Success = success
	resp.User = user

	return nil
}

func (s *ServiceContract) ListUsers(args *datacontract.EmptyArgs, resp *datacontract.UserList) error {
	resp.Users = ListUsers()

	return nil
}
