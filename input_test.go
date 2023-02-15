package main

import (
	"github.com/jfixby/pin"
	"github.com/jfixby/pin/fileops"
	"github.com/jfixby/vivcourt/api_input"
	"github.com/jfixby/vivcourt/input"
	"path/filepath"
	"testing"
)

/*
Read input data from file and print into console.
*/

func TestInput(t *testing.T) {
	home := fileops.Abs("")
	testData := filepath.Join(home, "test", "test3")
	testInput := filepath.Join(testData, "input", "in.json")
	//testOutput := filepath.Join(testData, "api_output", "expected.log")

	testListener := &InputTestListener{}
	input.ReadAll(testInput, testListener)

	//jsn, err := json.Marshal(allEvents)
	//if err != nil {
	//	panic(err)
	//}

	//pin.D("json", string(jsn))

	pin.D("EXIT")
}

type InputTestListener struct {
}

var allEvents []api_input.Event = []api_input.Event{}

func (t *InputTestListener) DoProcess(ev *api_input.Event) {

	pin.D("Event received", ev)

	allEvents = append(allEvents, *ev)

}

func (t *InputTestListener) Reset() {
}
