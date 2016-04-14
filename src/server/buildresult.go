package main

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
