package runtime

import (
	"github.com/MarvinJWendt/testza"
	"github.com/Vilsol/go-mlog/mlog"
	"github.com/rs/zerolog"
	"testing"
)

func init() {
	zerolog.SetGlobalLevel(zerolog.WarnLevel)
}

func TestTokenizer(t *testing.T) {
	result, _ := mlog.Tokenize("foo \"hello world\" bar # abc\n1234 lorem ipsum")

	testza.AssertEqual(t, []mlog.MLOGLine{
		{Instruction: []string{"foo", "\"hello world\"", "bar"}, Comment: " abc", SourceLine: 0},
		{Instruction: []string{"1234", "lorem", "ipsum"}, Comment: "", SourceLine: 1},
	}, result)
}
