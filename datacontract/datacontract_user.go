package datacontract

type UserType int

const (
	Admin UserType = iota
	Teacher
	Student
)
