package runtime

import (
	"github.com/MarvinJWendt/testza"
	"github.com/Vilsol/go-mlog/cli"
	"github.com/Vilsol/go-mlog/runtime"
	"github.com/rs/zerolog"
	"testing"
)

func init() {
	zerolog.SetGlobalLevel(zerolog.WarnLevel)
}

func TestMemory(t *testing.T) {
	objects := map[string]interface{}{
		"message1": &cli.Message{
			Name:   "message1",
			Output: nil,
		},
		"cell1": &cli.Memory{
			Name: "cell1",
			Size: 64,
			Data: make(map[int64]float64),
		},
	}

	operations, err := runtime.Parse(`write 1234 cell1 0
write 1.234 cell1 1
write "hello" cell1 2
write 100 cell1 1000
read m0 cell1 0
read m1 cell1 1
read m2 cell1 2
read m1000 cell1 1000`)

	if err != nil {
		panic(err)
	}

	context, counter := runtime.ConstructContext(objects)
	err = runtime.ExecuteContext(operations, context, counter)

	testza.AssertEqual(t, float64(1234), context.Variables["m0"].Value)
	testza.AssertEqual(t, 1.234, context.Variables["m1"].Value)
	testza.AssertEqual(t, float64(0), context.Variables["m2"].Value)
	testza.AssertEqual(t, float64(0), context.Variables["m1000"].Value)

	operations, err = runtime.Parse(`write 1234 cell2 0`)
	testza.AssertNoError(t, err)

	context, counter = runtime.ConstructContext(objects)

	testza.AssertPanics(t, func() {
		err = runtime.ExecuteContext(operations, context, counter)
		testza.AssertNoError(t, err)
	})

	operations, err = runtime.Parse(`write 1234 message1 0`)
	testza.AssertNoError(t, err)

	context, counter = runtime.ConstructContext(objects)

	testza.AssertPanics(t, func() {
		err = runtime.ExecuteContext(operations, context, counter)
		testza.AssertNoError(t, err)
	})

	operations, err = runtime.Parse(`read 1234 cell2 0`)
	testza.AssertNoError(t, err)

	context, counter = runtime.ConstructContext(objects)

	testza.AssertPanics(t, func() {
		err = runtime.ExecuteContext(operations, context, counter)
		testza.AssertNoError(t, err)
	})

	operations, err = runtime.Parse(`read 1234 message1 0`)
	testza.AssertNoError(t, err)

	context, counter = runtime.ConstructContext(objects)

	testza.AssertPanics(t, func() {
		err = runtime.ExecuteContext(operations, context, counter)
		testza.AssertNoError(t, err)
	})
}
