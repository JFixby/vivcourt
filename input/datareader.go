package input

import (
	"encoding/binary"
	"fmt"
	"github.com/jfixby/pin"
	"github.com/jfixby/vivcourt/orderbook"
	"io"
	"os"
)

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
		return nil, err
	}

	return data, err
}

func ParseEvent(data []byte) *orderbook.Event {
	result := &orderbook.Event{}

	pin.D("in", fmt.Sprintf("read %d bytes: %q", len(data), data))
	pin.D(" ")

	result.Symbol = orderbook.Symbol(data[1:4])

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
