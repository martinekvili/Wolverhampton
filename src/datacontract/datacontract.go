package datacontract

// EmptyArgs is the args for the Empty service
type EmptyArgs struct {
	Name   string
	Number int
}

// EmptyResult is the result type for the Empty service
type EmptyResult struct {
	Result string
}
