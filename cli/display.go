package cli

import (
	"encoding/json"
	"github.com/Vilsol/go-mlog/runtime"
	"github.com/disintegration/imaging"
	"github.com/fogleman/gg"
	"github.com/rs/zerolog/log"
	"image"
	"image/draw"
	"image/png"
	"math"
	"os"
	"path"
	"strconv"
	"strings"
)

func NewDisplay(config ObjectConfig) (*Display, error) {
	var options DisplayOptions
	if err := json.Unmarshal(config.Options, &options); err != nil {
		return nil, err
	}

	return &Display{
		Width:  options.Width,
		Height: options.Height,
		Scale:  options.Scale,
		Output: options.Output,
		SaveCurrentFrame: func(result image.Image, display *Display) {
			dir, f := path.Split(display.Output)
			split := strings.SplitN(f, ".", 2)
			name := split[0] + "-" + strconv.Itoa(display.FrameCount) + "." + split[1]
			file, err := os.OpenFile(path.Join(dir, name), os.O_RDWR|os.O_CREATE, 0755)
			if err != nil {
				log.Err(err).Msg("error saving image")
				return
			}

			if err := png.Encode(file, result); err != nil {
				log.Err(err).Msg("error saving image")
				return
			}

			log.Info().Int("frame", display.FrameCount).Str("path", file.Name()).Msg("saved")
		},
	}, nil
}

func (m *Display) DrawFlush(buffer []runtime.DrawStatement) {
	var img *image.RGBA

	if m.PreviousFrame != nil {
		source := m.PreviousFrame
		bounds := source.Bounds()
		img = image.NewRGBA(bounds)
		draw.Draw(img, bounds, source, bounds.Min, draw.Src)
	} else {
		img = image.NewRGBA(image.Rect(0, 0, m.Width, m.Height))
	}

	drawContext := gg.NewContextForRGBA(img)
	drawContext.SetLineCapSquare()

	for _, statement := range buffer {
		log.Trace().Interface("args", statement.Arguments).Interface("action", statement.Action).Msg("flushing")

		switch statement.Action {
		case runtime.DrawActionClear:
			r := toFloat64(statement.Arguments[0])
			g := toFloat64(statement.Arguments[1])
			b := toFloat64(statement.Arguments[2])
			drawContext.SetRGB255(int(r), int(g), int(b))
			drawContext.Clear()
			break
		case runtime.DrawActionColor:
			r := toFloat64(statement.Arguments[0])
			g := toFloat64(statement.Arguments[1])
			b := toFloat64(statement.Arguments[2])
			a := toFloat64(statement.Arguments[3])
			drawContext.SetRGBA255(int(r), int(g), int(b), int(a))
			break
		case runtime.DrawActionStroke:
			w := toFloat64(statement.Arguments[0])
			drawContext.SetLineWidth(w)
			break
		case runtime.DrawActionLine:
			x1 := toFloat64(statement.Arguments[0])
			y1 := toFloat64(statement.Arguments[1])
			x2 := toFloat64(statement.Arguments[2])
			y2 := toFloat64(statement.Arguments[3])
			drawContext.DrawLine(x1, y1, x2, y2)
			drawContext.Stroke()
			break
		case runtime.DrawActionRect:
			x := toFloat64(statement.Arguments[0])
			y := toFloat64(statement.Arguments[1])
			w := toFloat64(statement.Arguments[2])
			h := toFloat64(statement.Arguments[3])
			drawContext.DrawRectangle(x, y, w, h)
			drawContext.Fill()
			break
		case runtime.DrawActionLineRect:
			x := toFloat64(statement.Arguments[0])
			y := toFloat64(statement.Arguments[1])
			w := toFloat64(statement.Arguments[2])
			h := toFloat64(statement.Arguments[3])
			drawContext.DrawRectangle(x, y, w, h)
			drawContext.Stroke()
			break
		case runtime.DrawActionPoly:
			x := toFloat64(statement.Arguments[0])
			y := toFloat64(statement.Arguments[1])
			sides := toFloat64(statement.Arguments[2])
			radius := toFloat64(statement.Arguments[3])
			rotation := gg.Radians(toFloat64(statement.Arguments[4]) + 90)
			DrawRegularPolygon(drawContext, int(sides), x, y, radius, rotation)
			drawContext.Fill()
			break
		case runtime.DrawActionLinePoly:
			x := toFloat64(statement.Arguments[0])
			y := toFloat64(statement.Arguments[1])
			sides := toFloat64(statement.Arguments[2])
			radius := toFloat64(statement.Arguments[3])
			rotation := gg.Radians(toFloat64(statement.Arguments[4]) + 90)
			DrawRegularPolygon(drawContext, int(sides), x, y, radius, rotation)
			drawContext.Stroke()
			break
		case runtime.DrawActionTriangle:
			x1 := toFloat64(statement.Arguments[0])
			y1 := toFloat64(statement.Arguments[1])
			x2 := toFloat64(statement.Arguments[2])
			y2 := toFloat64(statement.Arguments[3])
			x3 := toFloat64(statement.Arguments[4])
			y3 := toFloat64(statement.Arguments[5])
			drawContext.MoveTo(x1, y1)
			drawContext.LineTo(x2, y2)
			drawContext.LineTo(x3, y3)
			drawContext.ClosePath()
			drawContext.Fill()
			break
		case runtime.DrawActionImage:
			// TODO Image drawing
			log.Warn().Msg("IMAGE")
			break
		default:
			panic("unknown draw action: " + statement.Action)
		}
	}

	var result image.Image = img
	if m.Scale != 1 {
		result = imaging.Resize(img, int(float64(m.Width)*m.Scale), int(float64(m.Height)*m.Scale), imaging.NearestNeighbor)
	}

	result = imaging.FlipV(result)

	m.SaveCurrentFrame(result, m)
	m.FrameCount++
	m.PreviousFrame = img
}

func toFloat64(n interface{}) float64 {
	switch cast := n.(type) {
	case float64:
		return cast
	case int64:
		return float64(cast)
	case string:
		float, err := strconv.ParseFloat(cast, 64)
		if err != nil {
			return 0
		}
		return float
	}

	return 0
}

func DrawRegularPolygon(ctx *gg.Context, n int, x, y, r, rotation float64) {
	angle := 2 * math.Pi / float64(n)
	rotation -= math.Pi / 2
	ctx.NewSubPath()
	for i := 0; i < n; i++ {
		a := rotation + angle*float64(i)
		ctx.LineTo(x+r*math.Cos(a), y+r*math.Sin(a))
	}
	ctx.ClosePath()
}
