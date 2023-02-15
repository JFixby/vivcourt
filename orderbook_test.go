package main

import (
	"github.com/jfixby/pin"
	"github.com/jfixby/pin/fileops"
	"github.com/jfixby/vivcourt/api_input"
	"github.com/jfixby/vivcourt/api_output"
	"github.com/jfixby/vivcourt/input"
	"github.com/jfixby/vivcourt/orderbook"
	"github.com/jfixby/vivcourt/output"
	"path/filepath"
	"testing"
)

var setup *testing.T

// Both component test and usage example
func TestOrderBook(t *testing.T) {
	setup = t

	// input data
	home := fileops.Abs("")
	testData := filepath.Join(home, "test", "test2")
	testInput := filepath.Join(testData, "input", "in.stream")
	testOutput := filepath.Join(testData, "api_output", "expected.log")

	// expected api_output
	// expected api_output
	expectedOutput := &output.Output{File: testOutput}
	expectedOutput.ReadAll()

	// TestEnvironment wraps and tests Book component
	testEnvironment := &TestEnvironment{
		expectedOutput: expectedOutput}

	//create book and subscribe it to TestEnvironment
	book := orderbook.NewBook()
	testEnvironment.book = book

	// expected input will be read as a file and converted into event stream
	// fed to test environment
	input.ReadAll(testInput, testEnvironment)

	book.Print()

	pin.D("EXIT")
}

type TestEnvironment struct {
	expectedOutput *output.Output
	book           *orderbook.Book
	counter        int
}

// Receives input events and feeds them to the Book
func (t *TestEnvironment) DoProcess(ev *api_input.Event) {
	pin.D("Input ", ev)
	t.book.DoUpdate(ev)
}

// Listents for events spawned by the Book and checks them against expected
func (t *TestEnvironment) OnBookEvent(e *api_output.BookSnapshot) {
	pin.D("Output", e)
	expectedEvent := t.expectedOutput.GetLevel(t.counter)

	check(setup, e, expectedEvent, t.counter)
	t.counter++
}

// compares expected api_output with actual
func check(
	setup *testing.T,
	actual *api_output.BookSnapshot,
	expected *api_output.BookSnapshot,
	counter int) {

	if !expected.Equal(actual) {

		pin.D(" counter", counter)
		pin.D("expected", expected)
		pin.D("  actual", actual)
		//setup.FailNow()
		panic("")
	}
}

// Resets book on each scenario
func (t *TestEnvironment) Reset() {
	t.counter = 0
	t.book.Reset()
}
