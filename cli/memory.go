package cli

import (
	"encoding/json"
)

func NewMemory(config ObjectConfig) (*Memory, error) {
	var options MemoryOptions
	if err := json.Unmarshal(config.Options, &options); err != nil {
		return nil, err
	}

	return &Memory{
		Name: config.Name,
		Size: options.Size,
		Data: make(map[int64]float64),
	}, nil
}

func (m *Memory) Write(value float64, position int64) {
	if position > m.Size-1 {
		return
	}

	m.Data[position] = value
}

func (m *Memory) Read(position int64) float64 {
	if val, ok := m.Data[position]; ok {
		return val
	}

	return 0
}
