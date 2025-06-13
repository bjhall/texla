package parser

type TokenKind int
const (
	Keyword TokenKind = iota
	Identifier
	Integer
	Float
	Comma
	Period
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
	MinusMinus
	Plus
	PlusPlus
	Mult
	Div
	Comment
	RightArrow
	Whitespace
	StringLiteral
	Range
	QuestionMark
	LogicAnd
	LogicOr
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
	case Period: return "Period"
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
	case MinusMinus: return "MinusMinus"
	case Plus: return "Plus"
	case PlusPlus: return "PlusPlus"
	case Mult: return "Mult"
	case Div: return "Div"
	case Comment: return "Comment"
	case RightArrow: return "RightArrow"
	case Whitespace: return "Whitespace"
	case StringLiteral: return "StringLiteral"
	case Range: return "Range"
	case QuestionMark: return "QuestionMark"
	case LogicAnd: return "LogicAnd"
	case LogicOr: return "LogicOr"
	case NoToken: return "NoToken"
	case Eof: return "Eof"

	default: return "???"
	}
}


type NodeType int
const (
	NoOpNodeType NodeType = iota
	NumNodeType
	BoolNodeType
	StringLiteralNodeType
	BinOpNodeType
	UnaryOpNodeType
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
	SliceLiteralNodeType
	SetLiteralNodeType
	ArgumentNodeType
	IndexedVarNodeType
	ForeachNodeType
	RangeNodeType
	FailNodeType
	ContinueNodeType
	BreakNodeType
	IncNodeType
	DecNodeType
)


func (s NodeType) String() string {
	switch s {
	case NoOpNodeType: return "NoOp"
	case NumNodeType: return "Num"
	case BoolNodeType: return "Bool"
	case StringLiteralNodeType: return "StringLiteral"
	case BinOpNodeType: return "BinOp"
	case UnaryOpNodeType: return "UnaryOp"
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
	case SliceLiteralNodeType: return "SliceLiteral"
	case SetLiteralNodeType: return "SetLiteral"
	case ArgumentNodeType: return "Argument"
	case IndexedVarNodeType: return "IndexedVar"
	case ForeachNodeType: return "Foreach"
	case RangeNodeType: return "Range"
	case FailNodeType: return "Fail"
	case ContinueNodeType: return "Continue"
	case BreakNodeType: return "Break"
	case IncNodeType: return "Inc"
	case DecNodeType: return "Dec"

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

