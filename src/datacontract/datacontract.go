package datacontract

// EmptyArgs is a more explanatory name for empty arguments
type EmptyArgs bool

type BuildResultArgs struct {
	JobID       int
	BuildResult bool
}

type CloseJobArgs struct {
	JobID int
}

type StartJobArgs struct {
	JobID    int
	FileName string
}

type JobStatusArgs struct {
	JobID       int
	JobNumInRow int
}

type RunResult int

const (
	Success RunResult = iota
	NotEnoughMemory
	NotEnoughTime
	Unknown
)

type RunResultArgs struct {
	JobID  int
	Result RunResult
}

type OutputMismatchLine struct {
	LineNumber int
	Expected   string
	Actual     string
}

type OutputMatchResultArgs struct {
	JobID      int
	Mismatches []OutputMismatchLine
}
