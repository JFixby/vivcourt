package input

import (
	"github.com/jfixby/pin"
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
		pin.E("failed to open file", err)
		r.runFlag = false
		panic(err)
	}

	for r.runFlag {
		data, err := ReadMessageData(file)

		if err != nil {
			pin.E("failed to read file", err)
			r.runFlag = false
			panic(err)
		}

		//pin.D("in", fmt.Sprintf("read %d bytes: %q", len(data), data))
		if len(data) == 0 {
			break
		}

		event := ParseEvent(data)
		if r.listener != nil {
			if event != nil {
				r.listener.DoProcess(event)
			}
		}

	}

	r.runFlag = false
}
