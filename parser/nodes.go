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
		return TypeInt{}
	case Float:
		return TypeFloat{}
	default:
		panic("UNREACHABLE")
	}
}

// Boolean literals
type BoolNode struct {
	Node
	token Token
	value bool
}

func (n *BoolNode) Print(level int) {
	indentation := strings.Repeat(" ", level*4)
	fmt.Println(indentation+ n.token.str)
}

func (n *BoolNode) Type() NodeType {
	return BoolNodeType
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

// Indexed variable node
type IndexedVarNode struct {
	Node
	token Token
	index Node
}

func (n *IndexedVarNode) Print(level int) {
	indentation := strings.Repeat(" ", level*4)
	fmt.Println(indentation + "Indexed Variable: " + n.token.str)
	fmt.Println(indentation + "Index:")
	n.index.Print(level+1)
}

func (n *IndexedVarNode) Type() NodeType {
	return IndexedVarNodeType
}

func (n *IndexedVarNode) Precedence() int {
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

// Unary operator node
type UnaryOpNode struct {
	Node
	op     Token
	expr   Node
}

func (n *UnaryOpNode) Print(level int) {
	indentation := strings.Repeat(" ", level*4)
	fmt.Println(indentation+ "UnaryOp", n.op.str)
	n.expr.Print(level+1)
}

func (n *UnaryOpNode) Type() NodeType {
	return UnaryOpNodeType
}

func (n *UnaryOpNode) Precedence() int {
	switch n.op.kind {
	case Not:
		return 6
	default:
		panic("Precedence not implemented for unary operator")
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

// Argument
type ArgumentNode struct {
	Node
	expr      Node
	paramName string
	order     int
	named     bool
	typ       Type
}

func (n *ArgumentNode) Print(level int) {
	indentation := strings.Repeat(" ", level*4)
	fmt.Println(indentation+ "Argument", "named:", n.named, "paramName:", n.paramName, "order:", n.order, "type:", n.typ)
	n.expr.Print(level+1)
}

func (n *ArgumentNode) Type() NodeType {
	return ArgumentNodeType
}

func (n *ArgumentNode) Precedence() int {
	return 7
}


// Function call
type FunctionCallNode struct {
	Node
	name          string
	arguments     []Node
	argumentOrder []int
}

func (n *FunctionCallNode) Print(level int) {
	indentation := strings.Repeat(" ", level*4)
	fmt.Println(indentation+ "FunctionCall", n.name)
	if len(n.arguments) == 0 {
		return
	}
	if len(n.argumentOrder) == 0 {
		fmt.Println(indentation+ "Unordered arguments:")
		for _, arg := range n.arguments {
			arg.Print(level+1)
		}
	} else {
		fmt.Println(indentation+ "Ordered arguments:")
		for _, idx := range n.argumentOrder {
			n.arguments[idx].Print(level+1)
		}
	}
}

func (n *FunctionCallNode) Type() NodeType {
	return FunctionCallNodeType
}

func (n *FunctionCallNode) Precedence() int {
	return 7
}

func (n *FunctionCallNode) appendArgumentOrder(idx int) {
	n.argumentOrder = append(n.argumentOrder, idx)
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
	name         string
	typ          Type
	hasDefault   bool
	defaultValue string
}

func (n *ParameterNode) Print(level int) {
	indentation := strings.Repeat(" ", level*4)
	fmt.Println(indentation+ "Parameter", n.name, n.typ, "HasDefault:", n.hasDefault, "Default value:", n.defaultValue)
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
	elseBody Node
	compType Type
}

func (n *IfNode) Print(level int) {
	indentation := strings.Repeat(" ", level*4)
	fmt.Println(indentation + "If")
	n.comp.Print(level+1)
	n.body.Print(level+1)
	if n.elseBody.Type() != NoOpNodeType {
		fmt.Println(indentation + "Else")
		n.elseBody.Print(level+1)
	}
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


// SLice literal node
type SliceLiteralNode struct {
	Node
	elements    []Node
	elementType Type
}

func (n *SliceLiteralNode) Print(level int) {
	indentation := strings.Repeat(" ", level*4)
	fmt.Println(indentation + "Slice literal,", n.elementType)
	for i, e := range n.elements {
		fmt.Println(indentation + "    " + "Element", i)
		e.Print(level+1)
	}
}

func (n *SliceLiteralNode) Type() NodeType {
	return SliceLiteralNodeType
}

func (n *SliceLiteralNode) Precedence() int {
	return 1000
}
