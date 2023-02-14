package orderbook

import (
	"fmt"
	"github.com/jfixby/pin"
	"github.com/jfixby/vivcourt/util"
)

type BookListener interface {
	OnBookEvent(*BookEvent)
}

type Book struct {
	BookListener BookListener
	markets      map[Symbol]*Market
	ordersById   map[OrderID]*OrderList
}

type Market struct {
	Symbol     Symbol
	buyOrders  *util.SkipList // price :-> order list, log N search
	sellOrders *util.SkipList // price :-> order list, log N search
}

func (b *Book) removeOrder(orderId OrderID) {
	//pin.D("remove", orderId)

	order, olist := b.findOrder(orderId)
	orderStack := olist.owner

	bestOrderID := findBestOrderID(orderStack, order.Side)
	_, bestOrderList := b.findOrder(bestOrderID)

	wasBestOrder := bestOrderList == olist

	olist.list.Delete(int(orderId))
	olist.totalQuantity = olist.totalQuantity - order.Quantity
	delete(b.ordersById, orderId)

	if olist.list.Len() == 0 {
		orderStack.Delete(key(order.Price))

	}

	if orderStack.Len() == 0 {
		//b.bestShallow(order.Side)
	} else {
		if wasBestOrder {
			//newBestOrderID := findBestOrderID(orderStack, order.Side)
			//o, l := b.findOrder(newBestOrderID)
			//	b.best(o, l.totalQuantity)
		}
	}

}

func (b *Book) findOrder(orderId OrderID) (order *Order, orderList *OrderList) {
	orderList = b.ordersById[orderId]
	if orderList == nil {
		pin.E("Missing order", orderId)
		pin.E("             ", b.ordersById)
		panic("")
	}
	v, _ := orderList.list.Get(int(orderId))
	order = v.(*Order)
	return
}

func findBestOrderID(orderStack *util.SkipList, side Side) OrderID {
	var best util.Iterator = nil
	if side == BUY {
		best = orderStack.SeekToLast()
	}
	if side == SELL {
		best = orderStack.SeekToFirst()
	}
	olist := best.Value().(*OrderList)

	return OrderID(olist.list.SeekToFirst().Key().(int))
}

type Order struct {
	OrderID  OrderID
	Quantity Quantity
	Price    Price
	Symbol   Symbol
	Side     Side
}

type OrderList struct {
	list          *util.SkipList // order id :-> order
	totalQuantity Quantity
	owner         *util.SkipList //
}

func (b *Book) DoUpdate(ev *Event) {
	if ev.OrderType == ADD {
		b.NewOrder(ev)
	}
	if ev.OrderType == UPDATE {
		b.UpdateOrder(ev)
	}
	if ev.OrderType == DELETE {
		b.CancelOrder(ev)
	}
	if ev.OrderType == EXECUTE {
		b.ExecuteOrder(ev)
	}
}

func (b *Book) NewOrder(ev *Event) {

	order := &Order{}
	order.OrderID = ev.OrderID
	order.Price = ev.Price
	order.Symbol = ev.Symbol
	order.Quantity = ev.Size
	order.Side = ev.Side

	b.append(order)

	//if b.orderIsTradeable(order) {
	//	b.execute(order)
	//} else {
	//b.append(order)
	//
	//}

}

func (b *Book) UpdateOrder(ev *Event) {
	orderOld, _ := b.findOrder(ev.OrderID)
	if orderOld == nil {
		panic(fmt.Sprintf("Order book is inconsistent. Order <%v> not found.", ev.OrderID))
	}

	b.removeOrder(ev.OrderID)

	order := &Order{}
	order.OrderID = ev.OrderID
	order.Price = ev.Price
	order.Symbol = ev.Symbol
	order.Quantity = ev.Size
	order.Side = ev.Side

	b.append(order)
}

func (b *Book) ExecuteOrder(ev *Event) {
	order, _ := b.findOrder(ev.OrderID)
	if order == nil {
		panic(fmt.Sprintf("Order book is inconsistent. Order <%v> not found.", ev.OrderID))
	}

	b.execute(order)

}

func (b *Book) CancelOrder(order *Event) {
	b.removeOrder(order.OrderID)
}

func (b *Book) onExecuteOrder(buy *Order, sell *Order, price Price, quantity Quantity) {
	bev := &BookEvent{}
	//bev.EventType = TRADE

	//bev.OrderIDBuy = buy.OrderID
	//
	//bev.OrderIDSell = sell.OrderID
	//
	//bev.Price = price
	//bev.Quantity = quantity

	b.BookListener.OnBookEvent(bev)
}

func (b *Book) getMarket(symbol Symbol) *Market {
	if b.markets == nil {
		b.markets = map[Symbol]*Market{}
		b.ordersById = map[OrderID]*OrderList{}
	}
	market := b.markets[symbol]
	if market == nil {
		market = &Market{Symbol: symbol}
		market.buyOrders = util.NewIntMap()
		market.sellOrders = util.NewIntMap()
		b.markets[symbol] = market
	}
	return market
}

func (b *Book) Reset() *Book {
	b.markets = nil
	b.ordersById = nil
	return b
}

func (b *Book) append(order *Order) {
	market := b.getMarket(order.Symbol)

	var orderStack *util.SkipList = nil

	if order.Side == BUY {
		orderStack = market.buyOrders
	}
	if order.Side == SELL {
		orderStack = market.sellOrders
	}

	list, ok := orderStack.Get(key(order.Price))
	if !ok {
		list = &OrderList{}
		orderStack.Set(key(order.Price), list)

		olist := list.(*OrderList)
		olist.owner = orderStack
		olist.list = util.NewIntMap()

	}

	olist := list.(*OrderList)
	olist.list.Set(int(order.OrderID), order)
	b.ordersById[order.OrderID] = olist

	olist.totalQuantity = olist.totalQuantity + order.Quantity

	var price Price = 0
	if order.Side == BUY {
		highestBid := b.highestBid(order.Symbol)
		price = unKey(highestBid.Key())
	}
	if order.Side == SELL {
		lowestAsk := b.lowestAsk(order.Symbol)
		price = unKey(lowestAsk.Key())
	}

	if price == order.Price {
		//b.best(order, olist.totalQuantity)
	}

	return
}

func (b *Book) execute(order *Order) {
	market := b.getMarket(order.Symbol)

	remainingQuantity := order.Quantity

	var orderStack *util.SkipList = nil
	var level util.Iterator = nil

	if order.Side == BUY {
		orderStack = market.sellOrders
	}
	if order.Side == SELL {
		orderStack = market.buyOrders
	}

	for orderStack.Len() > 0 {
		if remainingQuantity == 0 {
			break
		}

		if order.Side == BUY {
			level = orderStack.SeekToFirst()
		}
		if order.Side == SELL {
			level = orderStack.SeekToLast()
		}

		price := unKey(level.Key())
		if order.Side == BUY {
			if price > order.Price {
				break
			}
		}
		if order.Side == SELL {
			if price < order.Price {
				break
			}
		}
		orders := level.Value().(*OrderList)

		var buy *Order = nil
		var sell *Order = nil

		for orders.list.Len() > 0 {
			nextOrder := (orders.list.SeekToFirst().Value()).(*Order)

			if order.Side == BUY {
				buy = order
				sell = nextOrder
			}
			if order.Side == SELL {
				buy = nextOrder
				sell = order
			}

			if nextOrder.Quantity <= remainingQuantity {
				quantityToExecute := nextOrder.Quantity

				b.onExecuteOrder(buy, sell, price, quantityToExecute)
				remainingQuantity = remainingQuantity - quantityToExecute
				b.removeOrder(nextOrder.OrderID)

			} else {
				quantityToExecute := remainingQuantity

				b.onExecuteOrder(buy, sell, price, quantityToExecute)
				nextOrder.Quantity = nextOrder.Quantity - quantityToExecute
				remainingQuantity = remainingQuantity - quantityToExecute //should be 0
				orders.totalQuantity = orders.totalQuantity - quantityToExecute

				break
			}

		}
	}

	if remainingQuantity > 0 {
		order.Quantity = remainingQuantity

		b.append(order)
	}

}

func Invert(side Side) Side {
	if side == BUY {
		return SELL
	}
	if side == SELL {
		return BUY
	}
	panic("Invalid state")
}

func key(price Price) int {
	return int(price)
}

func unKey(i interface{}) Price {
	return Price(i.(int))
}

func (b *Book) orderIsTradeable(order *Order) bool {

	if order.Side == BUY {
		lowestAsk := b.lowestAsk(order.Symbol)
		if lowestAsk == nil {
			return false
		}
		if unKey(lowestAsk.Key()) <= order.Price {
			return true
		}
		return false
	}

	if order.Side == SELL {
		highestBid := b.highestBid(order.Symbol)
		if highestBid == nil {
			return false
		}
		if unKey(highestBid.Key()) >= order.Price {
			return true
		}
		return false
	}

	panic("Invalid state")
}

func (b *Book) highestBid(sim Symbol) util.Iterator {
	market := b.getMarket(sim)
	return market.buyOrders.SeekToLast()
}

func (b *Book) lowestAsk(sim Symbol) util.Iterator {
	market := b.getMarket(sim)
	return market.sellOrders.SeekToFirst()
}

func NewBook(l BookListener) *Book {
	return (&Book{BookListener: l}).Reset()
}
