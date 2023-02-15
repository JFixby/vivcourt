package orderbook

import (
	"fmt"
	"github.com/jfixby/pin"
	"github.com/jfixby/vivcourt/api_input"
	"github.com/jfixby/vivcourt/util"
)

type Book struct {
	markets    map[api_input.Symbol]*Market
	ordersById map[api_input.OrderID]*OrderList
}

func (b *Book) Print() {
	pin.D("OrderBook:")
	for k, v := range b.markets {
		pin.D("market", k)
		pin.D("      ", "buy")
		printOrders("      ", v.buyOrders)
		pin.D("      ", "sell")
		printOrders("      ", v.sellOrders)

	}

}

func printOrders(prefix string, list *util.SkipList) {

	if list.Len() == 0 {
		return
	}
	i := list.Iterator()
	i.Next()
	defer i.Close()
	for {
		k := i.Key()
		v := i.Value().(*OrderList)
		pin.D(prefix+" price", k)
		printOrderStack(prefix+prefix, v.list)
		if !i.Next() {
			break
		}
	}

}

func printOrderStack(prefix string, list *util.SkipList) {
	if list.Len() == 0 {
		return
	}
	i := list.Iterator()
	i.Next()
	defer i.Close()
	for {
		id := i.Key().(int)
		order := i.Value().(*Order)
		pin.D(prefix+fmt.Sprintf("ID<%v>", id), order.Quantity)
		if !i.Next() {
			break
		}
	}
}

type Market struct {
	Symbol     api_input.Symbol
	buyOrders  *util.SkipList // price :-> order list, log N search
	sellOrders *util.SkipList // price :-> order list, log N search
}

func (b *Book) removeOrder(orderId api_input.OrderID) {
	//pin.D("remove", orderId)

	order, olist := b.findOrder(orderId)
	if order == nil {
		return
	}
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

func (b *Book) findOrder(orderId api_input.OrderID) (order *Order, orderList *OrderList) {
	orderList = b.ordersById[orderId]
	if orderList == nil {
		//pin.E("Missing order", orderId)
		//pin.E("             ", b.ordersById)
		//panic("")
		return
	}
	v, _ := orderList.list.Get(int(orderId))
	order = v.(*Order)
	return
}

func findBestOrderID(orderStack *util.SkipList, side api_input.Side) api_input.OrderID {
	var best util.Iterator = nil
	if side == api_input.BUY {
		best = orderStack.SeekToLast()
	}
	if side == api_input.SELL {
		best = orderStack.SeekToFirst()
	}
	olist := best.Value().(*OrderList)

	return api_input.OrderID(olist.list.SeekToFirst().Key().(int))
}

type Order struct {
	OrderID  api_input.OrderID
	Quantity api_input.Quantity
	Price    api_input.Price
	Symbol   api_input.Symbol
	Side     api_input.Side
}

type OrderList struct {
	list          *util.SkipList // order id :-> order
	totalQuantity api_input.Quantity
	owner         *util.SkipList //
}

func (b *OrderList) String() string {
	return fmt.Sprintf("%v", b.totalQuantity)
}

func (b *Book) DoUpdate(ev *api_input.Event) {
	if ev.OrderType == api_input.ADD {
		b.AddOrder(ev)
	}
	if ev.OrderType == api_input.UPDATE {
		b.UpdateOrder(ev)
	}
	if ev.OrderType == api_input.DELETE {
		b.DeleteOrder(ev)
	}
	if ev.OrderType == api_input.EXECUTE {
		b.ExecuteOrder(ev)
	}
	//b.Print()
	//pin.D("")
}

func (b *Book) AddOrder(ev *api_input.Event) {
	if ev.Size == 0 {
		return
	}

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

func (b *Book) UpdateOrder(ev *api_input.Event) {
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

func (b *Book) ExecuteOrder(ev *api_input.Event) {
	order, _ := b.findOrder(ev.OrderID)
	if order == nil {
		panic(fmt.Sprintf("Order book is inconsistent. Order <%v> not found.", ev.OrderID))
	}

	order.Quantity = order.Quantity - ev.Size
	if order.Quantity == 0 {
		b.removeOrder(order.OrderID)
	}

}

func (b *Book) DeleteOrder(order *api_input.Event) {
	b.removeOrder(order.OrderID)
}

func (b *Book) onExecuteOrder(
	buy *Order,
	sell *Order,
	price api_input.Price,
	quantity api_input.Quantity) {

}

func (b *Book) getMarket(symbol api_input.Symbol) *Market {
	if b.markets == nil {
		b.markets = map[api_input.Symbol]*Market{}
		b.ordersById = map[api_input.OrderID]*OrderList{}
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

	if order.Side == api_input.BUY {
		orderStack = market.buyOrders
	}
	if order.Side == api_input.SELL {
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

	var price api_input.Price = 0
	if order.Side == api_input.BUY {
		highestBid := b.highestBid(order.Symbol)
		price = unKey(highestBid.Key())
	}
	if order.Side == api_input.SELL {
		lowestAsk := b.lowestAsk(order.Symbol)
		price = unKey(lowestAsk.Key())
	}

	if price == order.Price {
		//b.best(order, olist.totalQuantity)
	}

	return
}

//func (b *Book) execute(order *Order) {
//	market := b.getMarket(order.Symbol)
//
//	remainingQuantity := order.Quantity
//
//	var orderStack *util.SkipList = nil
//	var level util.Iterator = nil
//
//	if order.Side == api_input.BUY {
//		orderStack = market.sellOrders
//	}
//	if order.Side == api_input.SELL {
//		orderStack = market.buyOrders
//	}
//
//	for orderStack.Len() > 0 {
//		if remainingQuantity == 0 {
//			break
//		}
//
//		if order.Side == api_input.BUY {
//			level = orderStack.SeekToFirst()
//		}
//		if order.Side == api_input.SELL {
//			level = orderStack.SeekToLast()
//		}
//
//		price := unKey(level.Key())
//		if order.Side == api_input.BUY {
//			if price > order.Price {
//				break
//			}
//		}
//		if order.Side == api_input.SELL {
//			if price < order.Price {
//				break
//			}
//		}
//		orders := level.Value().(*OrderList)
//
//		var buy *Order = nil
//		var sell *Order = nil
//
//		for orders.list.Len() > 0 {
//			nextOrder := (orders.list.SeekToFirst().Value()).(*Order)
//
//			if order.Side == api_input.BUY {
//				buy = order
//				sell = nextOrder
//			}
//			if order.Side == api_input.SELL {
//				buy = nextOrder
//				sell = order
//			}
//
//			if nextOrder.Quantity <= remainingQuantity {
//				quantityToExecute := nextOrder.Quantity
//
//				b.onExecuteOrder(buy, sell, price, quantityToExecute)
//				remainingQuantity = remainingQuantity - quantityToExecute
//				b.removeOrder(nextOrder.OrderID)
//
//			} else {
//				quantityToExecute := remainingQuantity
//
//				b.onExecuteOrder(buy, sell, price, quantityToExecute)
//				nextOrder.Quantity = nextOrder.Quantity - quantityToExecute
//				remainingQuantity = remainingQuantity - quantityToExecute //should be 0
//				orders.totalQuantity = orders.totalQuantity - quantityToExecute
//
//				break
//			}
//
//		}
//	}
//
//	if remainingQuantity > 0 {
//		order.Quantity = remainingQuantity
//
//		b.append(order)
//	}
//
//}

func Invert(side api_input.Side) api_input.Side {
	if side == api_input.BUY {
		return api_input.SELL
	}
	if side == api_input.SELL {
		return api_input.BUY
	}
	panic("Invalid state")
}

func key(price api_input.Price) int {
	return int(price)
}

func unKey(i interface{}) api_input.Price {
	return api_input.Price(i.(int))
}

func (b *Book) orderIsTradeable(order *Order) bool {

	if order.Side == api_input.BUY {
		lowestAsk := b.lowestAsk(order.Symbol)
		if lowestAsk == nil {
			return false
		}
		if unKey(lowestAsk.Key()) <= order.Price {
			return true
		}
		return false
	}

	if order.Side == api_input.SELL {
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

func (b *Book) highestBid(sim api_input.Symbol) util.Iterator {
	market := b.getMarket(sim)
	return market.buyOrders.SeekToLast()
}

func (b *Book) lowestAsk(sim api_input.Symbol) util.Iterator {
	market := b.getMarket(sim)
	return market.sellOrders.SeekToFirst()
}

func NewBook() *Book {
	return (&Book{}).Reset()
}
