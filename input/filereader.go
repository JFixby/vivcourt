package input

import (
	"encoding/binary"
	"fmt"
	"github.com/jfixby/pin"
	"github.com/jfixby/vivcourt/orderbook"
	"io"
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

func ReadMessageData(file *os.File) ([]byte, error) {
	var Sequence uint32 = 0
	var Size uint32 = 0

	err := binary.Read(file, binary.LittleEndian, &Sequence)

	if err == io.EOF {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}
	err = binary.Read(file, binary.LittleEndian, &Size)
	if err != nil {
		return nil, err
	}

	data := make([]byte, Size)
	count, err := file.Read(data)
	if count == 0 {
		return nil, nil
	}
	return data, nil
}

func ParseEvent(data []byte) *orderbook.Event {
	result := &orderbook.Event{}

	pin.D("in", fmt.Sprintf("read %d bytes: %q", len(data), data))
	pin.D(" ")

	if data[0] == 'A' {
		result.OrderType = orderbook.ADD
	}
	if data[0] == 'E' {
		result.OrderType = orderbook.EXECUTE
	}
	if data[0] == 'U' {
		result.OrderType = orderbook.UPDATE
	}
	if data[0] == 'D' {
		result.OrderType = orderbook.DELETE
	}

	return result
}
