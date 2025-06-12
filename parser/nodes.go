package parser

import (
	"fmt"
	"strings"
)

type Node interface {
	Print(level int)
	Precedence() int
	Token() Token
}

type CommonNode struct {
	Node
	token Token
}

func (n *CommonNode) Token() Token {
	return n.token
}

// NoOp
type NoOpNode struct {
	Node
}

func (n *NoOpNode) Print(level int) {
	indentation := strings.Repeat(" ", level*4)
	fmt.Println(indentation + "NoOp")
}

func (n *NoOpNode) Precedence() int {
	return 0
}

func (n *NoOpNode) Token() Token {
	return Token{}
}

// Number literals
type NumNode struct {
	CommonNode
	token Token
}

func (n *NumNode) Print(level int) {
	indentation := strings.Repeat(" ", level*4)
	fmt.Println(indentation + n.token.str)
}

func (n *NumNode) Precedence() int {
	return 0
}

func (n *NumNode) NumType() Type {
	switch n.token.kind {
	case Integer:
		return TypeInt{}
	case Float:
		return TypeFloat{}
	default:
		panic("UNREACHABLE")
	}
}


// Boolean literals
type BoolNode struct {
	CommonNode
	token Token
}

func (n *BoolNode) Print(level int) {
	indentation := strings.Repeat(" ", level*4)
	fmt.Println(indentation + n.token.str)
}

func (n *BoolNode) Precedence() int {
	return 0
}

// String literals
type StringLiteralNode struct {
	Node
	token Token
}

func (n *StringLiteralNode) Print(level int) {
	indentation := strings.Repeat(" ", level*4)
	fmt.Println(indentation + "\"" + n.token.str + "\"")
}

func (n *StringLiteralNode) Precedence() int {
	return 0
}

// Variable node
type VarNode struct {
	CommonNode
	token Token
}

func (n *VarNode) Print(level int) {
	indentation := strings.Repeat(" ", level*4)
	fmt.Println(indentation + "Variable: " + n.token.str)
}

func (n *VarNode) Precedence() int {
	return 0
}

// Indexed variable node
type IndexedVarNode struct {
	CommonNode
	token Token
	index Node
}

func (n *IndexedVarNode) Print(level int) {
	indentation := strings.Repeat(" ", level*4)
	fmt.Println(indentation + "Indexed Variable: " + n.token.str)
	fmt.Println(indentation + "Index:")
	n.index.Print(level + 1)
}

func (n *IndexedVarNode) Precedence() int {
	return 0
}

// Binary operator node
type BinOpNode struct {
	CommonNode
	left  Node
	token Token
	right Node
}

func (n *BinOpNode) Print(level int) {
	indentation := strings.Repeat(" ", level*4)
	fmt.Println(indentation+"BinOp", n.token.str)
	n.left.Print(level + 1)
	n.right.Print(level + 1)
}

func (n *BinOpNode) Precedence() int {
	switch n.token.kind {
	case LogicAnd, LogicOr:
		return 1
	case Equal, NotEqual:
		return 2
	case Greater, Less, GreaterEqual, LessEqual:
		return 3
	case Plus, Minus:
		return 4
	case Mult, Div:
		return 5

	default:
		panic("Precedence not implemented for binary operator")
	}
}

// Unary operator node
type UnaryOpNode struct {
	CommonNode
	token Token
	expr  Node
}

func (n *UnaryOpNode) Print(level int) {
	indentation := strings.Repeat(" ", level*4)
	fmt.Println(indentation+"UnaryOp", n.token.str)
	n.expr.Print(level + 1)
}

func (n *UnaryOpNode) Precedence() int {
	switch n.token.kind {
	case Not:
		return 6
	default:
		panic("Precedence not implemented for unary operator")
	}
}

// Assignment node
type AssignNode struct {
	CommonNode
	left        Node
	token       Token
	right       Node
	declaration bool
	expression  bool
}

func (n *AssignNode) Print(level int) {
	indentation := strings.Repeat(" ", level*4)
	if n.declaration {
		fmt.Println(indentation + "VarDeclaration", "Expr:", n.expression)
	} else {
		fmt.Println(indentation + "VarAssignment", "Expr:", n.expression)
	}

	n.left.Print(level + 1)
	n.right.Print(level + 1)
}

func (n *AssignNode) Precedence() int {
	return 6
}

// Compound statements
type CompoundStatementNode struct {
	CommonNode
	token      Token
	children   []Node
	unusedVars []string
	scope      *Scope
}

func (n *CompoundStatementNode) Print(level int) {
	indentation := strings.Repeat(" ", level*4)
	fmt.Println(indentation + "CompoundStatement, (unused vars: " + strings.Join(n.unusedVars, ",") + ")" + "") // + "Return type: " + n.scope.returnType.String())

	for _, child := range n.children {
		child.Print(level + 1)
	}
}

func (n *CompoundStatementNode) Precedence() int {
	return 0
}

func (n *CompoundStatementNode) SetVarType(name string, typ Type) {
	n.scope.setSymbolType(name, typ)
}

// Function
type FunctionNode struct {
	CommonNode
	token      Token
	parameters Node
	body       Node
	returnType Type
	fallible   bool
}

func (n *FunctionNode) Print(level int) {
	indentation := strings.Repeat(" ", level*4)
	fmt.Println(indentation+"Function", n.token.str, "fallible?", n.fallible)
	n.parameters.Print(level + 1)
	n.body.Print(level + 1)
}

func (n *FunctionNode) Precedence() int {
	return 0
}

// Argument
type ArgumentNode struct {
	CommonNode
	token     Token
	expr      Node
	paramName string
	order     int
	named     bool
	typ       Type
}

func (n *ArgumentNode) Print(level int) {
	indentation := strings.Repeat(" ", level*4)
	fmt.Println(indentation+"Argument", "named:", n.named, "paramName:", n.paramName, "order:", n.order, "type:", n.typ)
	n.expr.Print(level + 1)
}

func (n *ArgumentNode) Precedence() int {
	return 7
}

// Function call
type FunctionCallNode struct {
	CommonNode
	token              Token
	name               string
	arguments          []Node
	isBuiltin          bool
	resolvedArgs       map[string]ArgumentNode
	resolvedReturnType Type
	errorHandled       bool
	generatorBody      Node
	generatorVar       VarNode
	generatorHasIdx    bool
	generatorIdxVar    VarNode
	errorBody          Node
}

func (n *FunctionCallNode) Print(level int) {
	indentation := strings.Repeat(" ", level*4)
	fmt.Println(indentation+"FunctionCall", n.name, "builtin?", n.isBuiltin, "error handled?", n.errorHandled)
	if len(n.arguments) == 0 && len(n.resolvedArgs) == 0 {
		return
	}
	if len(n.resolvedArgs) == 0 {
		fmt.Println(indentation + "Unresolved arguments:")
		for _, arg := range n.arguments {
			arg.Print(level + 1)
		}
	} else {
		fmt.Println(indentation + "Resolved arguments:")
		for argName, arg := range n.resolvedArgs {
			fmt.Println(indentation + "* " + argName)
			arg.Print(level + 1)
		}
	}
	
	if n.generatorBody != nil {
		fmt.Println(indentation + "Generator (var:"+n.generatorVar.token.str+"):")
		n.generatorBody.Print(level + 1)
	}

	if n.errorBody != nil {
		fmt.Println(indentation + "Error block:")
		n.errorBody.Print(level + 1)
	}
}

func (n *FunctionCallNode) Precedence() int {
	return 7
}

func (n *FunctionCallNode) setArgType(argName string, typ Type) {
	arg, found := n.resolvedArgs[argName]
	if !found {
		panic("Trying to set type of non-existing or unresolved function argument")
	}
	arg.typ = typ
	n.resolvedArgs[argName] = arg
}

func (n *FunctionCallNode) matchArgsToParams(parameters []ParameterNode) error {

	// Only do the resolution once
	if len(n.resolvedArgs) > 0 {
		return nil
	}
	paramArgs := make(map[string]ArgumentNode)
	for i, param := range parameters {
		found := false
		for j, a := range n.arguments {
			arg := a.(*ArgumentNode)
			if (arg.named && arg.paramName == param.name) || (!arg.named && arg.order == i) {
				paramArgs[param.name] = *(n.arguments[j].(*ArgumentNode))
				found = true
				break
			}
		}
		if !found {
			if !param.hasDefault {
				return fmt.Errorf("Value missing for argument %q (%s) of function %q", param.name, param.typ, n.name)
			}
			paramArgs[param.name] = param.CreateDefaultNode()
		}
	}

	n.resolvedArgs = paramArgs
	return nil
}

// Program node
type ProgramNode struct {
	Node
	functions []Node
	imports   map[string]bool
	preludes  map[string]bool
}

func (n *ProgramNode) Print(level int) {
	indentation := strings.Repeat(" ", level*4)
	fmt.Println(indentation + "Program")
	for _, arg := range n.functions {
		arg.Print(level + 1)
	}
}

func (n *ProgramNode) Precedence() int {
	return 10
}

func (n *ProgramNode) addImport(importName string) {
	n.imports[importName] = true
}

func (n *ProgramNode) Token() Token {
	return Token{}
}

// Parameter node
type ParameterNode struct {
	CommonNode
	token        Token
	name         string
	typ          Type
	hasDefault   bool
	defaultValue string
}

func (n *ParameterNode) Print(level int) {
	indentation := strings.Repeat(" ", level*4)
	fmt.Println(indentation+"Parameter", n.name, n.typ, "HasDefault:", n.hasDefault, "Default value:", n.defaultValue)
}

func (n *ParameterNode) Precedence() int {
	return 10
}

func (n *ParameterNode) CreateDefaultNode() ArgumentNode {
	switch n.typ.(type) {
	case TypeString:
		return ArgumentNode{expr: &StringLiteralNode{token: Token{kind: StringLiteral, str: n.defaultValue}}}
	case TypeInt:
		return ArgumentNode{expr: &NumNode{token: Token{kind: Integer, str: n.defaultValue}}}
	case TypeFloat:
		return ArgumentNode{expr: &NumNode{token: Token{kind: Float, str: n.defaultValue}}}
	case TypeBool:
		return ArgumentNode{expr: &BoolNode{token: Token{kind: Keyword, str: n.defaultValue}}}
	default:
		panic(fmt.Sprintf("Cannot construct default value node for parameter type %q\n", n.typ))
	}
}

// Parameter list node
type ParameterListNode struct {
	CommonNode
	token      Token
	parameters []ParameterNode
}

func (n *ParameterListNode) Print(level int) {
	indentation := strings.Repeat(" ", level*4)
	fmt.Println(indentation + "ParameterList")
	for _, arg := range n.parameters {
		arg.Print(level + 1)
	}
}

func (n *ParameterListNode) Precedence() int {
	return 10
}

// Return node
type ReturnNode struct {
	CommonNode
	token    Token
	expr     Node
	typ      Type
	function FunctionNode
}

func (n *ReturnNode) Print(level int) {
	indentation := strings.Repeat(" ", level*4)
	fmt.Println(indentation + "Return")
	n.expr.Print(level + 1)
}

func (n *ReturnNode) setType(typ Type) {
	n.typ = typ
}

func (n *ReturnNode) Precedence() int {
	return 100
}

// Continue node
type ContinueNode struct {
	CommonNode
	token Token
}

func (n *ContinueNode) Print(level int) {
	indentation := strings.Repeat(" ", level*4)
	fmt.Println(indentation + "Continue")
}

func (n *ContinueNode) Precedence() int {
	return 100
}

// Break node
type BreakNode struct {
	CommonNode
	token Token
}

func (n *BreakNode) Print(level int) {
	indentation := strings.Repeat(" ", level*4)
	fmt.Println(indentation + "Break")
}

func (n *BreakNode) Precedence() int {
	return 100
}

// Fail node
type FailNode struct {
	CommonNode
	token    Token
	expr     Node
	typ      Type
	function FunctionNode
}

func (n *FailNode) Print(level int) {
	indentation := strings.Repeat(" ", level*4)
	fmt.Println(indentation + "Fail")
	n.expr.Print(level + 1)
}

func (n *FailNode) setType(typ Type) {
	n.typ = typ
}

func (n *FailNode) Precedence() int {
	return 100
}

// If node
type IfNode struct {
	CommonNode
	token    Token
	comp     Node
	body     Node
	elseBody Node
	compType Type
}

func (n *IfNode) Print(level int) {
	indentation := strings.Repeat(" ", level*4)
	fmt.Println(indentation + "If")
	n.comp.Print(level + 1)
	n.body.Print(level + 1)
	_, noElse := n.elseBody.(*NoOpNode)
	if noElse {
		fmt.Println(indentation + "Else")
		n.elseBody.Print(level + 1)
	}
}

func (n *IfNode) setCompType(typ Type) {
	n.compType = typ
}

func (n *IfNode) Precedence() int {
	return 1000
}

// Foreach node
type ForeachNode struct {
	CommonNode
	token       Token
	iterator    Node
	variable    VarNode
	idxVariable VarNode
	hasIdx      bool
	body        Node
}

func (n *ForeachNode) Print(level int) {
	indentation := strings.Repeat(" ", level*4)
	fmt.Println(indentation + "Foreach, iterator:")
	n.iterator.Print(level + 1)
	fmt.Println(indentation + "Foreach, control variable:")
	n.variable.Print(level + 1)
	if n.hasIdx {
		fmt.Println(indentation + "Foreach, index variable:")
		n.idxVariable.Print(level + 1)
	}
	fmt.Println(indentation + "Foreach, body:")
	n.body.Print(level + 1)
}

func (n *ForeachNode) Precedence() int {
	return 1000
}

// SLice literal node
type SliceLiteralNode struct {
	CommonNode
	token      Token
	elements    []Node
	elementType Type
}

func (n *SliceLiteralNode) Print(level int) {
	indentation := strings.Repeat(" ", level*4)
	fmt.Println(indentation+"Slice literal,", n.elementType)
	for i, e := range n.elements {
		fmt.Println(indentation+"    "+"Element", i)
		e.Print(level + 1)
	}
}

func (n *SliceLiteralNode) Precedence() int {
	return 1000
}

// SLice literal node
type SetLiteralNode struct {
	CommonNode
	token      Token
	elements    []Node
	elementType Type
}

func (n *SetLiteralNode) Print(level int) {
	indentation := strings.Repeat(" ", level*4)
	fmt.Println(indentation+"Set literal,", n.elementType)
	for i, e := range n.elements {
		fmt.Println(indentation+"    "+"Element", i)
		e.Print(level + 1)
	}
}

func (n *SetLiteralNode) Precedence() int {
	return 1000
}

// Increment node
type IncNode struct {
	CommonNode
	varName string
	token Token
}

func (n *IncNode) Print(level int) {
	indentation := strings.Repeat(" ", level*4)
	fmt.Println(indentation+"Increment:", n.varName)
}

func (n *IncNode) Precedence() int {
	return 1000
}


// Decrement node
type DecNode struct {
	CommonNode
	varName string
	token   Token
}

func (n *DecNode) Print(level int) {
	indentation := strings.Repeat(" ", level*4)
	fmt.Println(indentation+"Decrement:", n.varName)
}

func (n *DecNode) Precedence() int {
	return 1000
}

// Range node
type RangeNode struct {
	Node
	token Token
	from  Node
	to    Node
	step  int
}

func (n *RangeNode) Print(level int) {
	indentation := strings.Repeat(" ", level*4)
	fmt.Println(indentation+"Range, step:", n.step)
	fmt.Println(indentation + "    " + "From:")
	n.from.Print(level + 1)
	fmt.Println(indentation + "    " + "To:")
	n.to.Print(level + 1)
}

func (n *RangeNode) Precedence() int {
	return 1000
}
