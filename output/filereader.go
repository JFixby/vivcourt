package output

import (
	"bufio"
	"github.com/jfixby/pin"
	"github.com/jfixby/vivcourt/api_output"
	"os"
)

type Output struct {
	File   string
	levels []*api_output.BookSnapshot
}

func (o *Output) ReadAll() error {

	pin.D("reading", o.File)
	file, err := os.Open(o.File)
	defer file.Close()
	if err != nil {
		return err
	}

	o.levels = []*api_output.BookSnapshot{}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		txt := scanner.Text()
		event, err := TryToParse(txt)
		if err != nil {
			pin.E("failed to read input", file)
			panic(err)
		}
		if event != nil {
			//pin.D("", event)
			o.levels = append(o.levels, event)
		}

	}
	return nil
}

func (o *Output) GetLevel(counter int) *api_output.BookSnapshot {
	list := o.levels
	if counter >= len(list) {
		return nil
	}
	return list[counter]
}

func (o *Output) Print() {
	pin.D("test data", o.levels)
}
