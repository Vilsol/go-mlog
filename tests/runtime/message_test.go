package runtime

import (
	"bytes"
	"github.com/MarvinJWendt/testza"
	"github.com/Vilsol/go-mlog/cli"
	"github.com/Vilsol/go-mlog/runtime"
	"github.com/rs/zerolog"
	"testing"
)

func init() {
	zerolog.SetGlobalLevel(zerolog.WarnLevel)
}

func TestMessage(t *testing.T) {
	buffer := bytes.NewBufferString("")

	objects := map[string]interface{}{
		"message1": &cli.Message{
			Name:   "message1",
			Output: buffer,
		},
		"cell1": &cli.Memory{},
	}

	operations, err := runtime.Parse(`print "hello"
print 1234
print 1.234
printflush message1`)
	testza.AssertNoError(t, err)

	context, counter := runtime.ConstructContext(objects)
	err = runtime.ExecuteContext(operations, context, counter)
	testza.AssertNoError(t, err)

	testza.AssertEqual(t, "hello12341.234", buffer.String())

	operations, err = runtime.Parse(`printflush message2`)
	testza.AssertNoError(t, err)

	context, counter = runtime.ConstructContext(objects)

	testza.AssertPanics(t, func() {
		err = runtime.ExecuteContext(operations, context, counter)
		testza.AssertNoError(t, err)
	})

	operations, err = runtime.Parse(`printflush cell1`)
	testza.AssertNoError(t, err)

	context, counter = runtime.ConstructContext(objects)

	testza.AssertPanics(t, func() {
		err = runtime.ExecuteContext(operations, context, counter)
		testza.AssertNoError(t, err)
	})
}
