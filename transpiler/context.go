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
)

type ContextBlock struct {
	Statements []MLOGStatement
	Extra      []MLOGStatement
}
