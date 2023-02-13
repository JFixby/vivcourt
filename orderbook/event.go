package orderbook

import "encoding/json"

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

type Executed struct {
	Price Price
	Size  Quantity
}

type BookEvent struct {
	SequenceNumber SequenceNumber
	Symbol         Symbol
	Buy            *[]Executed
	Sell           *[]Executed
}

func (e *BookEvent) String() string {
	b, _ := json.Marshal(e)
	return string(b)
}

func (e *BookEvent) Equal(a *BookEvent) bool {
	return e == a
}
