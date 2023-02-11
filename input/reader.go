package input

import "github.com/jfixby/vivcourt/orderbook"

type DataListener interface {
	Reset(scenario string)
	DoProcess(*orderbook.Event)
}

type DataReader interface {
	Subscribe(DataListener)
	Run()
	Stop()
}
