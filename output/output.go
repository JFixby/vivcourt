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

		event, err := TryToParse(txt)
		if err != nil {
			pin.E("failed to read input", file)
			panic(err)
		}
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
		//return &orderbook.BookEvent{EventType: orderbook.OVER}
		return nil
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

func TryToParse(txt string) (*orderbook.BookEvent, error) {
	if txt == "" {
		return nil, nil
	}

	if txt[0:1] == "#" {
		return nil, nil
	}

	arr := strings.Split(txt, ", ")

	result := &orderbook.BookEvent{}

	sequenceId, err := strconv.Atoi(arr[0])
	if err != nil {
		return nil, err
	}
	result.SequenceNumber = orderbook.SequenceNumber(sequenceId)

	result.Symbol = orderbook.Symbol(arr[1])

	f := strings.Index(txt, "[") + 1
	t := strings.Index(txt, "]")
	result.Buy, err = collectArray(string(txt[f:t]))
	if err != nil {
		return nil, err
	}

	txt = txt[t+1:]
	f = strings.Index(txt, "[") + 1
	t = strings.Index(txt, "]")
	result.Sell, err = collectArray(string(txt[f:t]))
	if err != nil {
		return nil, err
	}

	return result, nil

}

func collectArray(array string) (*[]orderbook.Executed, error) {
	result := []orderbook.Executed{}

	for {
		f := strings.Index(array, "(")
		if f == -1 {
			break
		}
		l := strings.Index(array, ")")
		if l == -1 {
			panic("")
		}
		fragment := array[f+1 : l]
		array = array[l+1:]

		xyStr := strings.Split(fragment, ", ")

		price, err := strconv.Atoi(xyStr[0])
		if err != nil {
			panic(err)
		}

		quantity, err := strconv.Atoi(xyStr[1])
		if err != nil {
			panic(err)
		}

		result = append(result, orderbook.Executed{
			Price: orderbook.Price(price),
			Size:  orderbook.Quantity(quantity),
		})

	}

	return &result, nil
}
