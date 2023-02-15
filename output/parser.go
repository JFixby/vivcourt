package output

import (
	"github.com/jfixby/vivcourt/api_output"
	"strconv"
	"strings"
)

func TryToParse(txt string) (*api_output.BookSnapshot, error) {
	if txt == "" {
		return nil, nil
	}

	if txt[0:1] == "#" {
		return nil, nil
	}

	arr := strings.Split(txt, ", ")

	result := &api_output.BookSnapshot{}

	levl, err := strconv.Atoi(arr[0])
	if err != nil {
		return nil, err
	}
	result.Level = api_output.Level(levl)

	result.Symbol = api_output.Symbol(arr[1])

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

func collectArray(array string) (*[]api_output.AskBid, error) {
	result := []api_output.AskBid{}

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

		result = append(result, api_output.AskBid{
			Price: api_output.Price(price),
			Size:  api_output.Quantity(quantity),
		})

	}

	return &result, nil
}
