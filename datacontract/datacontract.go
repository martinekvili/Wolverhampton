package datacontract

// EmptyArgs is a more explanatory name for empty arguments
type EmptyArgs bool

type BuildResultArgs struct {
	JobID  int
	Result BuildResult
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

type Error struct {
	Type string `xml:"type"`

	FileName     string `xml:"filename"`
	LineNumber   int    `xml:"linenumber"`
	ColumnNumber int    `xml:"columnnumber"`

	Code    string `xml:"code"`
	Message string `xml:"message"`
}

type BuildResult struct {
	Successful bool `xml:"successful"`

	ErrorList []Error `xml:"errorlist>error"`
}

type JobResult struct {
	JobID       int
	BuildInfo   BuildResult
	RunInfo     RunResult
	CompareInfo []OutputMismatchLine
}

type GetJobResultArgs struct {
	JobID int
}
