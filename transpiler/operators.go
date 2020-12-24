package transpiler

import "go/token"

var jumpOperators = map[token.Token]string{
	token.EQL: "equal",
	token.NEQ: "notEqual",
	token.LSS: "lessThan",
	token.LEQ: "lessThanEq",
	token.GTR: "greaterThan",
	token.GEQ: "greaterThanEq",
}

// TODO Convert to structs and a registry
var regularOperators = map[token.Token]string{
	token.ADD:        "add",
	token.ADD_ASSIGN: "add",
	token.SUB:        "sub",
	token.SUB_ASSIGN: "sub",
	token.MUL:        "mul",
	token.MUL_ASSIGN: "mul",
	token.QUO:        "div",
	token.QUO_ASSIGN: "div",
	// TODO //
	token.REM:        "mod",
	token.REM_ASSIGN: "mod",
	// TODO SQRT
	token.EQL:        "equal",
	token.NEQ:        "notEqual",
	token.LSS:        "lessThan",
	token.LEQ:        "lessThanEq",
	token.GTR:        "greaterThan",
	token.GEQ:        "greaterThanEq",
	token.LAND:       "land",
	token.SHL:        "shl",
	token.SHL_ASSIGN: "shl",
	token.SHR:        "shr",
	token.SHR_ASSIGN: "shr",
	token.LOR:        "or",
	token.OR:         "or",
	token.OR_ASSIGN:  "or",
	token.AND:        "and",
	token.AND_ASSIGN: "and",
	token.XOR:        "xor",
	token.XOR_ASSIGN: "xor",
	token.NOT:        "not",
	// TODO max
	// TODO min
	// TODO atan2
	// TODO dst
	// TODO noise
	// TODO abs
	// TODO log
	// TODO log10
	// TODO sin
	// TODO cos
	// TODO tan
	// TODO floor
	// TODO ceil
	// TODO sqrt
	// TODO rand
}
