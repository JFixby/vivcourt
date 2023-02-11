package testoutput

import (
	"bufio"
	"github.com/jfixby/pin"
	"github.com/jfixby/vivcourt/orderbook"
	"os"
	"strconv"
	"strings"
)

type TestOutput struct {
	File string
	data map[string][]*orderbook.BookEvent
}

func (o *TestOutput) LoadAll() error {
	o.data = map[string][]*orderbook.BookEvent{}

	pin.D("reading", o.File)
	file, err := os.Open(o.File)
	defer file.Close()
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(file)
	tag := ""
	for scanner.Scan() {
		txt := scanner.Text()

		if strings.HasPrefix(txt, "#name: ") {
			tag = txt[len("#name: "):]
			{
				//pin.D("tag", tag)
				o.data[tag] = []*orderbook.BookEvent{}
				continue
			}
		}

		event := TryToParse(txt)
		if event != nil {
			//pin.D("", event)
			o.data[tag] = append(o.data[tag], event)
		}

	}
	return nil
}

func (o *TestOutput) GetEvent(scenario string, counter int) *orderbook.BookEvent {
	list := o.data[scenario]
	if list == nil {
		pin.E("scenario not found", scenario)
		pin.E("                  ", counter)
		pin.E("                  ", o.data)
		panic("")
	}
	if counter == len(list) {
		return &orderbook.BookEvent{EventType: orderbook.OVER}
	}
	if counter >= len(list) {
		pin.E("output not found  ", o.data)
		return nil
	}
	return list[counter]
}

func (o *TestOutput) Print() {
	for k, v := range o.data {
		pin.D("test data", k)
		pin.D("", v)
	}
}

func TryToParse(txt string) *orderbook.BookEvent {
	if txt == "" {
		return nil
	}

	if txt[0:1] == "#" {
		return nil
	}

	arr := strings.Split(txt, ", ")

	result := &orderbook.BookEvent{}

	if arr[0] == "A" {
		result.EventType = orderbook.ACKNOWLEDGE

		userID, err := strconv.Atoi(arr[1])
		if err != nil {
			pin.E("invalid input", txt)
			return nil
		}
		result.UserIDAcknowledge = orderbook.UserID(userID)

		orderID, err := strconv.Atoi(arr[2])
		if err != nil {
			pin.E("invalid input", txt)
			return nil
		}
		result.OrderIDAcknowledge = orderbook.OrderID(orderID)
		return result
	}

	if arr[0] == "B" {
		result.EventType = orderbook.BEST

		if arr[1] == "S" {
			result.Side = orderbook.SELL
		} else if arr[1] == "B" {
			result.Side = orderbook.BUY
		} else {
			pin.E("Unknown order side", txt)
			return nil
		}

		if arr[2] == "-" {
			result.ShallowAsk = true
		} else {
			price, err := strconv.Atoi(arr[2])
			if err != nil {
				pin.E("invalid input", txt)
				return nil
			}
			result.Price = orderbook.Price(price)
		}

		if arr[3] == "-" {
			result.ShallowAsk = true
		} else {
			quantity, err := strconv.Atoi(arr[3])
			if err != nil {
				pin.E("invalid input", txt)
				return nil
			}
			result.Quantity = orderbook.Quantity(quantity)
		}

		return result
	}

	if arr[0] == "R" {
		result.EventType = orderbook.REJECT

		userID, err := strconv.Atoi(arr[1])
		if err != nil {
			pin.E("invalid input", txt)
			return nil
		}
		result.UserIDReject = orderbook.UserID(userID)

		orderID, err := strconv.Atoi(arr[2])
		if err != nil {
			pin.E("invalid input", txt)
			return nil
		}
		result.OrderIDReject = orderbook.OrderID(orderID)
		return result
	}

	if arr[0] == "T" {
		result.EventType = orderbook.TRADE

		userIDBuy, err := strconv.Atoi(arr[1])
		if err != nil {
			pin.E("invalid input", txt)
			return nil
		}
		result.UserIDBuy = orderbook.UserID(userIDBuy)

		orderIDBuy, err := strconv.Atoi(arr[2])
		if err != nil {
			pin.E("invalid input", txt)
			return nil
		}
		result.OrderIDBuy = orderbook.OrderID(orderIDBuy)

		userIDSell, err := strconv.Atoi(arr[3])
		if err != nil {
			pin.E("invalid input", txt)
			return nil
		}
		result.UserIDSell = orderbook.UserID(userIDSell)

		orderIDSell, err := strconv.Atoi(arr[4])
		if err != nil {
			pin.E("invalid input", txt)
			return nil
		}
		result.OrderIDSell = orderbook.OrderID(orderIDSell)

		price, err := strconv.Atoi(arr[5])
		if err != nil {
			pin.E("invalid input", txt)
			return nil
		}
		result.Price = orderbook.Price(price)

		quantity, err := strconv.Atoi(arr[6])
		if err != nil {
			pin.E("invalid input", txt)
			return nil
		}
		result.Quantity = orderbook.Quantity(quantity)
		return result
	}

	return nil
}
