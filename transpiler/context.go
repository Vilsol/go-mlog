package transpiler

type ContextKey string

const (
	contextOptions           = ContextKey("options")
	contextGlobal            = ContextKey("global")
	contextFunction          = ContextKey("function")
	contextStatement         = ContextKey("statement")
	contextSpec              = ContextKey("spec")
	contextDecl              = ContextKey("decl")
	contextBlock             = ContextKey("block")
	contextBreakableBlock    = ContextKey("breakableBlock")
	contextSwitchClauseBlock = ContextKey("switchClauseBlock")
	typeError                = ContextKey("typeError")
)

type ContextBlock struct {
	Statements []MLOGStatement
	Extra      []MLOGStatement
}
