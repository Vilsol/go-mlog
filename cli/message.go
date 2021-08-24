package cli

import (
	"encoding/json"
	"fmt"
	"os"
)

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
		Name:   config.Name,
		Output: output,
	}, nil
}

func (m Message) PrintFlush(buffer string) {
	if m.Output != nil {
		_, _ = m.Output.WriteString(fmt.Sprintf("[%s] %s\n", m.Name, buffer))
	}
}
