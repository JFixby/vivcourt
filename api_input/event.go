package api_input

import "encoding/json"

type EventConsumer interface {
	Reset()
	DoProcess(*Event)
}

type OrderType string

const ADD OrderType = "ADD"
const UPDATE OrderType = "UPDATE"
const DELETE OrderType = "DELETE"
const EXECUTE OrderType = "EXECUTE"

type Side string

const BUY Side = "BUY"
const SELL Side = "SELL"

type Symbol string
type OrderID int64
type Price int64
type Quantity int64
type SequenceNumber uint32

type Event struct {
	SequenceNumber SequenceNumber
	OrderType      OrderType
	Symbol         Symbol
	Price          Price
	Size           Quantity
	Side           Side
	OrderID        OrderID
}

func (e *Event) String() string {
	b, _ := json.Marshal(e)
	return string(b)
}
