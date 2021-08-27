package runtime

import (
	"github.com/Vilsol/go-mlog/cli"
	"github.com/Vilsol/go-mlog/runtime"
	"github.com/rs/zerolog"
	"os"
	"testing"
)

func init() {
	zerolog.SetGlobalLevel(zerolog.WarnLevel)
}

func BenchmarkFactorial(b *testing.B) {
	message1 := &cli.Message{
		Name:   "message1",
		Output: nil,
	}

	message2 := &cli.Message{
		Name:   "message1",
		Output: nil,
	}

	objects := map[string]interface{}{
		"message1": message1,
		"message2": message2,
	}

	input, _ := os.ReadFile("../../samples/factorial.mlog")
	operations, err := runtime.Parse(string(input))

	if err != nil {
		panic(err)
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		context, counter := runtime.ConstructContext(objects)
		if err := runtime.ExecuteContext(operations, context, counter); err != nil {
			panic(err)
		}
	}
}
