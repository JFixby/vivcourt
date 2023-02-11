package input

import (
	"github.com/jfixby/pin"
	"github.com/jfixby/vivcourt/orderbook"
	"os"
)

type FileReader struct {
	inputFile string
	listener  DataListener
	runFlag   bool
}

func NewFileReader(inputFile string) *FileReader {
	return &FileReader{
		inputFile,
		nil,
		false,
	}
}

func (r *FileReader) Subscribe(l DataListener) {
	r.listener = l
}

func (r *FileReader) IsRunnung() bool {
	return r.runFlag
}

func (r *FileReader) Stop() {
	r.runFlag = false
}

func (r *FileReader) Run() {
	if r.runFlag {
		return
	}

	r.runFlag = true
	go r.runthread()
}

func (r *FileReader) runthread() {
	input := r.inputFile
	pin.D("reading", input)
	file, err := os.Open(input)
	defer file.Close()
	if err != nil {
		pin.E("failed to read file", err)
		r.runFlag = false
		panic(err)
	}

	data := make([]byte, 1024)

	for r.runFlag {

		count, err := file.Read(data)
		if err != nil {
			pin.E("failed to read file", err)
			r.runFlag = false
			panic(err)
		}

		if count == 0 {
			break
		}

		pin.D("read %d bytes: %q\n", count, data[:count])

		event := ParseEvent(data)
		if r.listener != nil {
			if event != nil {
				r.listener.DoProcess(event)
			}
		}

	}

	r.runFlag = false
}

func ParseEvent(dt []byte) *orderbook.Event {
	result := &orderbook.Event{}
	return result
}
