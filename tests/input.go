package main

import (
	"github.com/h2non/survey"
	"github.com/h2non/survey/tests/util"
)

var val = ""

var table = []TestUtil.TestTableEntry{
	{
		"no default", &survey.Input{Message: "Hello world"}, &val,
	},
	{
		"default", &survey.Input{Message: "Hello world", Default: "default"}, &val,
	},
	{
		"no help, send '?'", &survey.Input{Message: "Hello world"}, &val,
	},
	{
		"input text in random location", &survey.Input{Message: "Hello"}, &val,
	},
}

func main() {
	TestUtil.RunTable(table)
}
