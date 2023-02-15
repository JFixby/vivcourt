package input

import (
	"encoding/binary"
	"github.com/jfixby/vivcourt/api_input"
	"io"
)

func ReadMessageData(r io.Reader) (Sequence uint32, data []byte, err error) {
	//var Sequence uint32 = 0
	var Size uint32 = 0

	err = binary.Read(r, binary.LittleEndian, &Sequence)
	if err == io.EOF {
		return
	}
	if err != nil {
		return
	}

	err = binary.Read(r, binary.LittleEndian, &Size)
	if err != nil {
		return
	}

	data = make([]byte, Size)
	count, err := r.Read(data)
	if count == 0 {
		return
	}

	return
}

func ParseEvent(sequence uint32, data []byte) *api_input.Event {
	result := &api_input.Event{}

	//	pin.D("in", fmt.Sprintf("read %d bytes: %q", len(data), string(data)))
	//pin.D("  ", fmt.Sprintf("               %v",hex.EncodeToString(data[4:4+8])))

	result.SequenceNumber = api_input.SequenceNumber(sequence)

	result.Symbol = api_input.Symbol(data[1:4])

	if data[0] == 'A' {
		result.OrderType = api_input.ADD
		result.Price = api_input.Price(binary.LittleEndian.Uint32(data[24 : 24+4]))
	}
	if data[0] == 'E' {
		result.OrderType = api_input.EXECUTE

		result.Size = api_input.Quantity(binary.LittleEndian.Uint64(data[16 : 16+8]))
	}
	if data[0] == 'U' {
		result.OrderType = api_input.UPDATE
		result.Price = api_input.Price(binary.LittleEndian.Uint32(data[24 : 24+4]))
	}
	if data[0] == 'D' {
		result.OrderType = api_input.DELETE
	}

	if result.OrderType != api_input.DELETE {
		result.Size = api_input.Quantity(binary.LittleEndian.Uint64(data[16 : 16+8]))
	}

	if data[12] == 'B' {
		result.Side = api_input.BUY
	}
	if data[12] == 'S' {
		result.Side = api_input.SELL
	}

	result.OrderID = api_input.OrderID(binary.LittleEndian.Uint64(data[4 : 4+8]))

	return result
}
