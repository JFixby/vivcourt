package input

import (
	"encoding/binary"
	"github.com/jfixby/vivcourt/orderbook"
	"io"
	"os"
)

func ReadMessageData(file *os.File) (Sequence uint32, data []byte, err error) {
	//var Sequence uint32 = 0
	var Size uint32 = 0

	err = binary.Read(file, binary.LittleEndian, &Sequence)
	if err == io.EOF {
		return
	}
	if err != nil {
		return
	}

	err = binary.Read(file, binary.LittleEndian, &Size)
	if err != nil {
		return
	}

	data = make([]byte, Size)
	count, err := file.Read(data)
	if count == 0 {
		return
	}

	return
}

func ParseEvent(sequence uint32, data []byte) *orderbook.Event {
	result := &orderbook.Event{}

	//	pin.D("in", fmt.Sprintf("read %d bytes: %q", len(data), string(data)))
	//pin.D("  ", fmt.Sprintf("               %v",hex.EncodeToString(data[4:4+8])))

	result.SequenceNumber = orderbook.SequenceNumber(sequence)

	result.Symbol = orderbook.Symbol(data[1:4])

	if data[0] == 'A' {
		result.OrderType = orderbook.ADD
		result.Price = orderbook.Price(binary.LittleEndian.Uint32(data[24 : 24+4]))
	}
	if data[0] == 'E' {
		result.OrderType = orderbook.EXECUTE

		result.Size = orderbook.Quantity(binary.LittleEndian.Uint64(data[16 : 16+8]))
	}
	if data[0] == 'U' {
		result.OrderType = orderbook.UPDATE
		result.Price = orderbook.Price(binary.LittleEndian.Uint32(data[24 : 24+4]))
	}
	if data[0] == 'D' {
		result.OrderType = orderbook.DELETE
	}

	if result.OrderType != orderbook.DELETE {
		result.Size = orderbook.Quantity(binary.LittleEndian.Uint64(data[16 : 16+8]))
	}

	if data[12] == 'B' {
		result.Side = orderbook.BUY
	}
	if data[12] == 'S' {
		result.Side = orderbook.SELL
	}

	result.OrderID = orderbook.OrderID(binary.LittleEndian.Uint64(data[4 : 4+8]))

	return result
}
