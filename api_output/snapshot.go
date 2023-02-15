package api_output

import (
	"encoding/json"
)

type Symbol string
type Price int64
type Quantity int64
type Level int

type AskBid struct {
	Price Price
	Size  Quantity
}

type BookSnapshot struct {
	Level  Level
	Symbol Symbol
	Buy    *[]AskBid
	Sell   *[]AskBid
}

func (e *BookSnapshot) String() string {
	b, _ := json.Marshal(e)
	return string(b)
}

func (e *BookSnapshot) Equal(a *BookSnapshot) bool {
	return e == a
}
