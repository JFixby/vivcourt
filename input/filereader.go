package input

import (
	"bufio"
	"encoding/json"
	"github.com/jfixby/pin"
	"github.com/jfixby/vivcourt/api_input"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

func ReadAll(input string, listener api_input.EventConsumer) {
	pin.D("reading", input)

	if input == "" {
		reader := bufio.NewReader(os.Stdin)
		ParseBinary(reader, listener)
		return
	}

	file, err := os.Open(input)
	defer file.Close()
	if err != nil {
		pin.E("failed to open file", err)
		panic(err)
	}
	if strings.HasSuffix(input, ".json") {
		ParseJson(file, listener)
	} else if strings.HasSuffix(input, ".stream") {
		ParseBinary(file, listener)
	}

}

func ParseJson(r io.Reader, listener api_input.EventConsumer) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		panic(err)
	}

	allEvents := []api_input.Event{}

	err = json.Unmarshal(b, &allEvents)
	if err != nil {
		panic(err)
	}

	for _, event := range allEvents {
		if listener != nil {
			listener.DoProcess(&event)
		}
	}

}

func ParseBinary(r io.Reader, listener api_input.EventConsumer) {
	for {
		sequence, data, err := ReadMessageData(r)

		if err == io.EOF {
			break
		}

		if err != nil {
			pin.E("failed to read file", err)
			panic(err)
		}

		//pin.D("in", fmt.Sprintf("read %d bytes: %q", len(data), data))
		if len(data) == 0 {
			break
		}

		event := ParseEvent(sequence, data)
		if listener != nil {
			if event != nil {
				listener.DoProcess(event)
			}
		}

	}
}
