package cli

import (
	"encoding/json"
	"github.com/Vilsol/go-mlog/runtime"
	"image"
	"io"
)

// Compile-time checks
var (
	_ runtime.Message = (*Message)(nil)
	_ runtime.Display = (*Display)(nil)
	_ runtime.Memory  = (*Memory)(nil)
)

type ObjectType string

const (
	ObjectMessage = ObjectType("message")
	ObjectDisplay = ObjectType("display")
	ObjectMemory  = ObjectType("memory")
)

type Config struct {
	Objects []ObjectConfig
}

type ObjectConfig struct {
	Type    ObjectType
	Name    string
	Options json.RawMessage
}

type MessageOptions struct {
	Output string
}

type Message struct {
	Name   string
	Output io.Writer
}

type DisplayOptions struct {
	Width  int
	Height int
	Scale  float64
	Output string
}

type Display struct {
	Width            int
	Height           int
	Scale            float64
	Output           string
	FrameCount       int
	PreviousFrame    image.Image
	SaveCurrentFrame func(img image.Image, display *Display)
}

type MemoryOptions struct {
	Size int64
}

type Memory struct {
	Name string
	Size int64
	Data map[int64]float64
}
