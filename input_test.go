package main

import (
	"github.com/jfixby/pin"
	"github.com/jfixby/pin/fileops"
	"github.com/jfixby/vivcourt/input"
	"github.com/jfixby/vivcourt/orderbook"
	"path/filepath"
	"testing"
	"time"
)

/*
Read input data from file and print into console.
*/

func TestInput(t *testing.T) {
	home := fileops.Abs("")
	testData := filepath.Join(home, "test", "test1")
	testInput := filepath.Join(testData, "input", "in.stream")
	//testOutput := filepath.Join(testData, "output", "expected.log")

	reader := input.NewFileReader(testInput)
	testListener := &InputTestListener{}
	reader.Subscribe(testListener)
	reader.Run()

	for reader.IsRunnung() {
		time.Sleep(2 * time.Second)
	}

	pin.D("EXIT")
}

type InputTestListener struct {
}

func (t *InputTestListener) DoProcess(ev *orderbook.Event) {
	pin.D("Event received", ev)
	pin.D(" ")
}

func (t *InputTestListener) Reset(scenario string) {
	pin.D("Next scenario", scenario)

}
