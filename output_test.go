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

	test := &output.Output{File: testOutput}

	err := test.ReadAll()

	if err != nil {
		t.Fatal(err)
	}

	test.Print()

	pin.D("EXIT")
}
