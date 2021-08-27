package runtime

import (
	"github.com/MarvinJWendt/testza"
	"github.com/Vilsol/go-mlog/cli"
	"github.com/Vilsol/go-mlog/runtime"
	"github.com/disintegration/imaging"
	"github.com/rs/zerolog"
	"image"
	"image/png"
	"os"
	"testing"
)

func init() {
	zerolog.SetGlobalLevel(zerolog.WarnLevel)
}

func TestDisplay(t *testing.T) {
	frames := make([]image.Image, 0)
	objects := map[string]interface{}{
		"display1": &cli.Display{
			Width:  176,
			Height: 176,
			Scale:  1,
			SaveCurrentFrame: func(result image.Image, display *cli.Display) {
				frames = append(frames, result)
			},
		},
		"cell1": &cli.Memory{},
	}

	operations, err := runtime.Parse(`draw clear 0 128 255
draw color 255 128 0 255
draw triangle 5 10 20 50 30 20
draw color 0 128 0 255
draw rect 45 10 25 25
draw color 0 255 128 255
draw stroke 1
draw lineRect 95 10 25 25
draw color 255 0 0 128
draw rect 140 10 25 25
draw stroke 8
draw line 45 45 165 50
draw color 255 255 0 255
draw poly 20 75 3 13 0
draw poly 55 75 4 13 0
draw poly 90 75 5 13 0
draw poly 125 75 6 13 0
draw poly 160 75 7 13 0
draw color 255 0 128 255
draw stroke 5
draw linePoly 20 125 3 13 0
draw linePoly 55 125 4 13 0
draw linePoly 90 125 5 13 0
draw linePoly 125 125 6 13 0
draw linePoly 160 125 7 13 0
draw color 255 0 0 255
draw poly 10 160 3 7 0
draw poly 25 160 3 7 24
draw poly 40 160 3 7 48
draw poly 55 160 3 7 72
draw poly 70 160 3 7 96
draw poly 85 160 3 7 130
draw poly 100 160 3 7 154
draw poly 115 160 3 7 178
draw poly 130 160 3 7 202
draw poly 145 160 3 7 226
draw poly 160 160 3 7 248
drawflush display1`)
	testza.AssertNoError(t, err)

	context, counter := runtime.ConstructContext(objects)
	err = runtime.ExecuteContext(operations, context, counter)
	testza.AssertNoError(t, err)

	testza.AssertEqual(t, 1, len(frames))

	expectedDisplay1File, err := os.Open("data/display_1.png")
	testza.AssertNoError(t, err)

	expectedDisplay1, err := png.Decode(expectedDisplay1File)
	testza.AssertNoError(t, err)

	expectedDisplay1NRGBA := imaging.Clone(expectedDisplay1)

	for x := 0; x < expectedDisplay1NRGBA.Bounds().Dx(); x++ {
		for y := 0; y < expectedDisplay1NRGBA.Bounds().Dx(); y++ {
			testza.AssertEqual(t, expectedDisplay1NRGBA.At(x, y), frames[0].At(x, y))
		}
	}

	testza.AssertPanics(t, func() {
		operations, err := runtime.Parse(`drawflush display2`)
		testza.AssertNoError(t, err)
		context, counter := runtime.ConstructContext(objects)
		err = runtime.ExecuteContext(operations, context, counter)
		testza.AssertNoError(t, err)
	})

	testza.AssertPanics(t, func() {
		operations, err := runtime.Parse(`drawflush cell1`)
		testza.AssertNoError(t, err)
		context, counter := runtime.ConstructContext(objects)
		err = runtime.ExecuteContext(operations, context, counter)
		testza.AssertNoError(t, err)
	})
}
