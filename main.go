package main

import (
	"github.com/jfixby/pin"
	"github.com/jfixby/vivcourt/api_input"
	"github.com/jfixby/vivcourt/input"
	"github.com/jfixby/vivcourt/orderbook"
)

func main() {

	book := orderbook.NewBook()
	setup := &Setup{book: book}
	input.ReadAll("", setup)
}

type Setup struct {
	book    *orderbook.Book
	counter int
}

func (t *Setup) DoProcess(ev *api_input.Event) {
	pin.D("Input ", ev)
	t.book.DoUpdate(ev)
}

func (t *Setup) Reset() {
	t.counter = 0
	t.book.Reset()
}
