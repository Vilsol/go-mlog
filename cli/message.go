package cli

import (
	"encoding/json"
	"io"
	"os"
)

type MessageHijacker struct {
	Prefix string
	Output io.Writer
}

func (h MessageHijacker) Write(p []byte) (n int, err error) {
	return h.Output.Write(append(append([]byte("["+h.Prefix+"] "), p...), byte('\n')))
}

func NewMessage(config ObjectConfig) (*Message, error) {
	output := os.Stdout

	var options MessageOptions
	if err := json.Unmarshal(config.Options, &options); err != nil {
		return nil, err
	}

	if options.Output == "stderr" {
		output = os.Stderr
	}

	return &Message{
		Name: config.Name,
		Output: &MessageHijacker{
			Prefix: config.Name,
			Output: output,
		},
	}, nil
}

func (m Message) PrintFlush(buffer string) {
	if m.Output != nil {
		_, _ = m.Output.Write([]byte(buffer))
	}
}
