package parser

import (
	"fmt"
	"strings"
)

type Node interface {
	Print(level int)
	Type() NodeType
	Precedence() int
}

// NoOp 
type NoOpNode struct {
	Node
}

func (n *NoOpNode) Print(level int) {
	indentation := strings.Repeat(" ", level*4)
	fmt.Println(indentation+ "NoOp")
}

func (n *NoOpNode) Type() NodeType {
	return NoOpNodeType
}

func (n *NoOpNode) Precedence() int {
	return 0
}


// Number literals
type NumNode struct {
	Node
	token Token
}

func (n *NumNode) Print(level int) {
	indentation := strings.Repeat(" ", level*4)
	fmt.Println(indentation+ n.token.str)
}

func (n *NumNode) Type() NodeType {
	return NumNodeType
}

func (n *NumNode) Precedence() int {
	return 0
}

func (n *NumNode) NumType() Type {
	switch n.token.kind {
	case Integer:
		return TypeInt
	case Float:
		return TypeFloat
	default:
		panic("UNREACHABLE")
	}
}

// String literals
type StringLiteralNode struct {
	Node
	token Token
}

func (n *StringLiteralNode) Print(level int) {
	indentation := strings.Repeat(" ", level*4)
	fmt.Println(indentation + "\""+ n.token.str + "\"")
}

func (n *StringLiteralNode) Type() NodeType {
	return StringLiteralNodeType
}

func (n *StringLiteralNode) Precedence() int {
	return 0
}


// Variable node
type VarNode struct {
	Node
	token Token
}

func (n *VarNode) Print(level int) {
	indentation := strings.Repeat(" ", level*4)
	fmt.Println(indentation + "Variable: " + n.token.str)
}

func (n *VarNode) Type() NodeType {
	return VarNodeType
}

func (n *VarNode) Precedence() int {
	return 0
}


// Binary operator node
type BinOpNode struct {
	Node
	left   Node
	op     Token
	right  Node
}

func (n *BinOpNode) Print(level int) {
	indentation := strings.Repeat(" ", level*4)
	fmt.Println(indentation+ "BinOp", n.op.str)
	n.left.Print(level+1)
	n.right.Print(level+1)
}

func (n *BinOpNode) Type() NodeType {
	return BinOpNodeType
}

func (n *BinOpNode) Precedence() int {
	switch n.op.kind {
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


// Assignment node
type AssignNode struct {
	Node
	left        Node
	tok         Token
	right       Node
	declaration bool
}

func (n *AssignNode) Print(level int) {
	indentation := strings.Repeat(" ", level*4)
	if n.declaration {
		fmt.Println(indentation+ "VarDeclaration")
	} else {
		fmt.Println(indentation+ "VarAssignment")
	}
	n.left.Print(level+1)
	n.right.Print(level+1)
}

func (n *AssignNode) Type() NodeType {
	return AssignNodeType
}

func (n *AssignNode) Precedence() int {
	return 6
}


// Compound statements
type CompoundStatementNode struct {
	Node
	children   []Node
	unusedVars []string
	scope      *Scope
}

func (n *CompoundStatementNode) Print(level int) {
	indentation := strings.Repeat(" ", level*4)
	fmt.Println(indentation+ "CompoundStatement, (unused vars: " + strings.Join(n.unusedVars, ",") +")" + "")// + "Return type: " + n.scope.returnType.String())

	for _, child := range n.children {
		child.Print(level+1)
	}
}

func (n *CompoundStatementNode) Type() NodeType {
	return CompoundStatementNodeType
}

func (n *CompoundStatementNode) Precedence() int {
	return 0
}

// Function
type FunctionNode struct {
	Node
	name       Token
	parameters Node
	body       Node
	returnType Type
}

func (n *FunctionNode) Print(level int) {
	indentation := strings.Repeat(" ", level*4)
	fmt.Println(indentation+ "Function", n.name.str)
	n.parameters.Print(level+1)
	n.body.Print(level+1)
}

func (n *FunctionNode) Type() NodeType {
	return FunctionNodeType
}

func (n *FunctionNode) Precedence() int {
	return 0
}


// Function call
type FunctionCallNode struct {
	Node
	name          string
	arguments     []Node
	argumentTypes []Type
}

func (n *FunctionCallNode) Print(level int) {
	indentation := strings.Repeat(" ", level*4)
	fmt.Println(indentation+ "FunctionCall", n.name)
	for _, arg := range n.arguments {
		arg.Print(level+1)
	}
	for i, typ := range n.argumentTypes {
		fmt.Println(indentation+"    Type:",i,typ)
	}
}

func (n *FunctionCallNode) Type() NodeType {
	return FunctionCallNodeType
}

func (n *FunctionCallNode) Precedence() int {
	return 7
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
	fmt.Println(indentation+ "Program")
	for _, arg := range n.functions {
		arg.Print(level+1)
	}
}

func (n *ProgramNode) Type() NodeType {
	return ProgramNodeType
}

func (n *ProgramNode) Precedence() int {
	return 10
}

func (n *ProgramNode) addImport(importName string) {
	n.imports[importName] = true
}

// Parameter node
type ParameterNode struct {
	Node
	name string
	typ  Type
}

func (n *ParameterNode) Print(level int) {
	indentation := strings.Repeat(" ", level*4)
	fmt.Println(indentation+ "Parameter", n.name, n.typ)
}

func (n *ParameterNode) Type() NodeType {
	return ParameterNodeType
}

func (n *ParameterNode) Precedence() int {
	return 10
}

// Parameter list node
type ParameterListNode struct {
	Node
	parameters []Node
}

func (n *ParameterListNode) Print(level int) {
	indentation := strings.Repeat(" ", level*4)
	fmt.Println(indentation+ "ParameterList")
	for _, arg := range n.parameters {
		arg.Print(level+1)
	}
}

func (n *ParameterListNode) Type() NodeType {
	return ParameterListNodeType
}

func (n *ParameterListNode) Precedence() int {
	return 10
}


// Return node
type ReturnNode struct {
	Node
	expr     Node
	typ      Type
	function FunctionNode
}

func (n *ReturnNode) Print(level int) {
	indentation := strings.Repeat(" ", level*4)
	fmt.Println(indentation+ "Return")
	n.expr.Print(level+1)
}

func (n *ReturnNode) setType(typ Type) {
	n.typ = typ
}

func (n *ReturnNode) Type() NodeType {
	return ReturnNodeType
}

func (n *ReturnNode) Precedence() int {
	return 100
}


// If node
type IfNode struct {
	Node
	comp     Node
	body     Node
	compType Type
}

func (n *IfNode) Print(level int) {
	indentation := strings.Repeat(" ", level*4)
	fmt.Println(indentation + "If")
	n.comp.Print(level+1)
	n.body.Print(level+1)
}

func (n *IfNode) setCompType(typ Type) {
	n.compType = typ
}

func (n *IfNode) Type() NodeType {
	return IfNodeType
}

func (n *IfNode) Precedence() int {
	return 1000
}
