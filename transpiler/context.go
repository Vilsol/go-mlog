package transpiler

const (
	contextOptions           = "options"
	contextGlobal            = "global"
	contextFunction          = "function"
	contextStatement         = "statement"
	contextSpec              = "spec"
	contextDecl              = "decl"
	contextBlock             = "block"
	contextBreakableBlock    = "breakableBlock"
	contextSwitchClauseBlock = "switchClauseBlock"
	typeError                = "typeError"
)

type ContextBlock struct {
	Statements []MLOGStatement
	Extra      []MLOGStatement
}
