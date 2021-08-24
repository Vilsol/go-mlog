package cli

import (
	"encoding/json"
	"github.com/Vilsol/go-mlog/runtime"
	"image"
	"os"
)

// Compile-time checks
var _ runtime.Message = (*Message)(nil)
var _ runtime.Display = (*Display)(nil)

type ObjectType string

const (
	ObjectMessage = ObjectType("message")
	ObjectDisplay = ObjectType("display")
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
	Output *os.File
}

type DisplayOptions struct {
	Width  int
	Height int
	Scale  float64
	Output string
}

type Display struct {
	Width         int
	Height        int
	Scale         float64
	Output        string
	FrameCount    int
	PreviousFrame image.Image
}
