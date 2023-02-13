package main

import (
	"github.com/jfixby/pin"
	"github.com/jfixby/pin/fileops"
	"github.com/jfixby/vivcourt/output"
	"path/filepath"
	"testing"
)

func TestOutput(t *testing.T) {
	home := fileops.Abs("")
	testData := filepath.Join(home, "test", "test1")
	testOutput := filepath.Join(testData, "output", "expected.log")

	test := &testoutput.TestOutput{File: testOutput}

	test.LoadAll()

	test.Print()

	pin.D("EXIT")
}
