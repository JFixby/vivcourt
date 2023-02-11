package orderbook

import "encoding/json"

type OrderType string

const NEW OrderType = "NEW"
const CANCEL OrderType = "CANCEL"
const FLUSH OrderType = "FLUSH"

type Side string

const BUY Side = "BUY"
const SELL Side = "SELL"

type Symbol string
type UserID int64
type OrderID int64
type Price int64
type Quantity int64

type Event struct {
	OrderType OrderType
	UserID    UserID
	Symbol    Symbol
	Price     Price
	Quantity  Quantity
	Side      Side
	OrderID   OrderID
}

func (e *Event) String() string {
	b, _ := json.Marshal(e)
	return string(b)
}

type BookEvent struct {
	//Input *Event

	EventType EventType

	UserIDAcknowledge UserID
	UserIDSell        UserID
	UserIDBuy         UserID
	UserIDReject      UserID

	Price    Price
	Quantity Quantity
	Side     Side

	OrderIDBuy         OrderID
	OrderIDSell        OrderID
	OrderIDAcknowledge OrderID
	OrderIDReject      OrderID
	ShallowAsk         bool
}

func (e *BookEvent) String() string {
	b, _ := json.Marshal(e)
	return string(b)
}

func (e *BookEvent) Equal(a *BookEvent) bool {
	return *e == *a
}

type EventType string

const ACKNOWLEDGE EventType = "ACKNOWLEDGE"
const REJECT EventType = "REJECT"
const BEST EventType = "BEST"
const TRADE EventType = "TRADE"
const OVER EventType = "OVER"
