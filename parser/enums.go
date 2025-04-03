package parser

type TokenKind int
const (
	Keyword TokenKind = iota
	Identifier
	Integer
	Float
	Comma
	OpenCurly
	CloseCurly
	OpenParen
	CloseParen
	OpenBracket
	CloseBracket
	DoubleQuote
	Greater
	GreaterEqual
	Less
	LessEqual
	Equal
	Not
	NotEqual
	Assign
	Minus
	Plus
	Mult
	Div
	Comment
	RightArrow
	Whitespace
	StringLiteral
	NoToken
	Eof
)


func (s TokenKind) String() string {
	switch s {
	case Keyword: return "Keyword"
	case Identifier: return "Identifier"
	case Integer: return "Integer"
	case Float: return "Float"
	case Comma: return "Comma"
	case OpenCurly: return "OpenCurly"
	case CloseCurly: return "CloseCurly"
	case OpenParen: return "OpenParen"
	case CloseParen: return "CloseParen"
	case OpenBracket: return "OpenBracket"
	case CloseBracket: return "CloseBracket"
	case DoubleQuote: return "DoubleQuote"
	case Greater: return "Greater"
	case GreaterEqual: return "GreaterEqual"
	case Less: return "Less"
	case LessEqual: return "LessEqual"
	case Equal: return "Equal"
	case Not: return "Not"
	case NotEqual: return "NotEqual"
	case Assign: return "Assign"
	case Minus: return "Minus"
	case Plus: return "Plus"
	case Mult: return "Mult"
	case Div: return "Div"
	case Comment: return "Comment"
	case RightArrow: return "RightArrow"
	case Whitespace: return "Whitespace"
	case StringLiteral: return "StringLiteral"
	case NoToken: return "NoToken"
	case Eof: return "Eof"

	default: return "???"
	}
}


type NodeType int
const (
	NoOpNodeType NodeType = iota
	NumNodeType
	StringLiteralNodeType
	BinOpNodeType
	CompoundStatementNodeType
	StatementNodeType
	AssignNodeType
	VarNodeType
	FunctionNodeType
	FunctionCallNodeType
	ProgramNodeType
	ParameterNodeType
	ParameterListNodeType
	ReturnNodeType
	IfNodeType
)


func (s NodeType) String() string {
	switch s {
	case NoOpNodeType: return "NoOp"
	case NumNodeType: return "Num"
	case StringLiteralNodeType: return "StringLiteral"
	case BinOpNodeType: return "BinOp"
	case CompoundStatementNodeType: return "CompoundStatement"
	case StatementNodeType: return "Statement"
	case AssignNodeType: return "Assign"
	case VarNodeType: return "Var"
	case FunctionNodeType: return "Function"
	case FunctionCallNodeType: return "FunctionCall"
	case ProgramNodeType: return "Program"
	case ParameterNodeType: return "Parameter"
	case ParameterListNodeType: return "ParameterList"
	case ReturnNodeType: return "Return"
	case IfNodeType: return "If"

	default: return "???"
	}
}


type Type int
const (
	NoCoercion Type = iota
	NoReturnType
	TypeUndetermined
	TypeInt
	TypeFloat
	TypeString
	TypeBool
)


func (s Type) String() string {
	switch s {
	case NoCoercion: return "NoCoercion"
	case NoReturnType: return "NoReturn"
	case TypeUndetermined: return "Undetermined"
	case TypeInt: return "Int"
	case TypeFloat: return "Float"
	case TypeString: return "String"
	case TypeBool: return "Bool"

	default: return "???"
	}
}


type CoercionMode int
const (
	CoercionModeDefault CoercionMode = iota
	CoercionModeNumLiteral
)
func (s CoercionMode) String() string {
	switch s {
	case CoercionModeDefault: return "Default"
	case CoercionModeNumLiteral: return "NumLiteral"

	default: return "???"
	}
}

