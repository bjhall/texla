package parser

import (
	"fmt"
)

type Symbol struct {
	typ            Type
	name           string
	used           bool
	category       SymbolCategory
	parameterTypes []Type
}

func (v *Symbol) setUsed() {
	v.used = true
}

func (v *Symbol) getType(typ Type) Type {
	return v.typ
}



type Scope struct {
	parent     *Scope
	symbols    map[string]Symbol
	returnType Type
}

type SymbolCategory int
const (
	VariableSymbol SymbolCategory = iota
	FunctionSymbol
)

func newScope(parent *Scope, parameters *ParameterListNode, returnType Type) *Scope {

	symbols := make(map[string]Symbol)

	// Add function parameters to the scopes list of declared symbols
	if parameters != nil {
		for _, param := range parameters.parameters {
			name := param.(*ParameterNode).name
			typ := param.(*ParameterNode).typ
			symbols[name] = Symbol{typ, name, false, VariableSymbol, []Type{}}
		}
	}

	return &Scope{parent, symbols, returnType}
}

func (s *Scope) setSymbolType(name string, typ Type) error {
	if symbol, exists := s.symbols[name]; exists {
		symbol.typ = typ
		s.symbols[name] = symbol
		return nil
	}
	return fmt.Errorf("UNREACHABLE: Trying to set type of non-existing symbol %q", name)
}

func (s *Scope) lookupSymbol(name string) (Symbol, bool) {
	if symbol, exists := s.symbols[name]; exists {
		return symbol, true
	}
	if s.parent == nil {
		return Symbol{}, false
	}
	return s.parent.lookupSymbol(name)
}

func (s *Scope) createSymbol(name string, category SymbolCategory, typ Type, parameterTypes []Type) bool {
	if _, exists := s.symbols[name]; exists {
		return false
	}
	s.symbols[name] = Symbol{typ, name, false, category, parameterTypes} // FIXME: Don't hardcode Type
	return true
}




type Parser struct {
	tokens       []Token
	tokenIdx     int
	blockDepth   int
	currentScope *Scope
	imports      map[string]bool
}

func (p *Parser) addImport(name string) {
	p.imports[name] = true
}

func (s *Scope) setVariableUsed(name string) error {
	if symbol, exists := s.symbols[name]; exists {
		symbol.used = true
		s.symbols[name] = symbol
		return nil
	}
	if s.parent == nil {
		return fmt.Errorf("Variable not found: %q", name)
	}
	return s.parent.setVariableUsed(name)
}


func (p *Parser) validateVariable(name string) bool {
	symbol, found := p.currentScope.lookupSymbol(name)

	// Symbol not found or was not a variable
	if !found || symbol.category != VariableSymbol {
		return false
	}

	// Annotate the variable as used
	p.currentScope.setVariableUsed(name)

	return true
}


func (p *Parser) createVariableInCurrentScope(name string, typ Type) bool {
	return p.currentScope.createSymbol(name, VariableSymbol, typ, []Type{})
}

func (p *Parser) createFunctionInCurrentScope(name string, parameterList *ParameterListNode, returnType Type) bool {
	var parameterTypes []Type
	for _, param := range parameterList.parameters {
		parameterTypes = append(parameterTypes, param.(*ParameterNode).typ)
	}
	return p.currentScope.createSymbol(name, FunctionSymbol, returnType, parameterTypes)
}

func (p *Parser) unusedVariables() []string {
	var unused []string
	for name, symbolData := range p.currentScope.symbols {
		if symbolData.category == VariableSymbol && !symbolData.used {
			unused = append(unused, name)
		}
	}
	return unused
}

func (p *Parser) newScope(parameters *ParameterListNode, returnType Type) {
	p.currentScope = newScope(p.currentScope, parameters, returnType)
}

func (p *Parser) leaveScope() {
	p.currentScope = p.currentScope.parent
}

func (p *Parser) currentToken() Token {
	return p.tokens[p.tokenIdx]
}

func (p *Parser) peek(i int) Token {
	return p.tokens[p.tokenIdx+i]
}

func (p *Parser) consumeToken() Token {
	p.tokenIdx++
	return p.tokens[p.tokenIdx-1]
}

func (p *Parser) expectToken(kind TokenKind) (Token, error) {
	if p.currentToken().kind == kind {
		return p.consumeToken(), nil
	}
	return Token{}, fmt.Errorf("Invalid token: expected %s, got %s %d:%d", kind, p.currentToken().kind, p.currentToken().line, p.currentToken().column)
}



func (p *Parser) parseExpr() (Node, error) {
	node, err := p.parseComparison()
	if err != nil {
		return &NoOpNode{}, err
	}
	return node, nil
}

func (p *Parser) parseComparison() (Node, error) {
	node, err := p.parseTerm()
	if err != nil {
		return &NoOpNode{}, err
	}

	for p.currentToken().kind == Equal || p.currentToken().kind == NotEqual || p.currentToken().kind == Greater || p.currentToken().kind == GreaterEqual ||  p.currentToken().kind == Less ||  p.currentToken().kind == LessEqual {
		opToken := p.consumeToken()
		right, err := p.parseTerm()
		if err != nil {
			return &NoOpNode{}, err
		}
		node = &BinOpNode{left: node, op: opToken, right: right}
	}
	return node, nil
}

func (p *Parser) parseTerm() (Node, error) {
	node, err := p.parseFactor()
	if err != nil {
		return &NoOpNode{}, err
	}
	for p.currentToken().kind == Plus || p.currentToken().kind == Minus {
		opToken := p.consumeToken()
		right, err := p.parseFactor()
		if err != nil {
			return &NoOpNode{}, err
		}
		node = &BinOpNode{left: node, op: opToken, right: right}
	}
	return node, nil
}


func (p *Parser) parseFactor() (Node, error) {
	node, err := p.parsePrimary()
	if err != nil {
		return &NoOpNode{}, err
	}

	for p.currentToken().kind == Mult || p.currentToken().kind == Div {
		opToken := p.consumeToken()
		right, err := p.parsePrimary()
		if err != nil {
			return &NoOpNode{}, err
		}
		node = &BinOpNode{left: node, op: opToken, right: right}
	}
	return node, nil
}

func (p *Parser) parsePrimary() (Node, error) {
	switch p.currentToken().kind {

	case OpenParen:
		_ = p.consumeToken()
		expr, err := p.parseExpr()
		if err != nil {
			return &NumNode{}, err
		}
		_ = p.consumeToken()
		return expr, nil

	case Integer, Float:
		return &NumNode{token: p.consumeToken()}, nil

	case Identifier:
		switch p.peek(1).kind {
		case OpenParen: // Function call
			functionCall, err := p.parseFunctionCall()
			if err != nil {
				return &NoOpNode{}, err
			}
			return functionCall, nil
		default: // Variable
			variable, err := p.parseVar(true)
			if err != nil {
				return &NoOpNode{}, err
			}
			return variable, nil
		}

	case StringLiteral:
		return &StringLiteralNode{token: p.consumeToken()}, nil

	default:
		panic("Invalid token, expected Num")
	}
}


func (p *Parser) parseCompoundStatement(parameters *ParameterListNode, returnType Type) (Node, error) {

	_, err := p.expectToken(OpenCurly)
	if err != nil {
		return &NoOpNode{}, err
	}

	// Create new scope
	p.newScope(parameters, returnType)
	defer p.leaveScope()

	// Parse any statements
	var statements []Node
	for p.currentToken().kind != CloseCurly {
		node, err := p.parseStatement()
		if err != nil {
			return &NoOpNode{}, err
		}
		statements = append(statements, node)
	}

	_, err = p.expectToken(CloseCurly)
	if err != nil {
		return &NoOpNode{}, err
	}
	return &CompoundStatementNode{children: statements, unusedVars: p.unusedVariables(), scope: p.currentScope}, nil
}

func (p *Parser) parseType() (Type, error) {
	typeToken, err := p.expectToken(Identifier)
	if err != nil {
		return TypeUndetermined, err
	}
	switch typeToken.str {
	case "int":
		return TypeInt, nil
	case "float":
		return TypeFloat, nil
	case "str":
		return TypeString, nil
	default:
		return TypeUndetermined, fmt.Errorf("Unknown type: %q", typeToken.str)
	}
}

func (p *Parser) parseParameter() (Node, error) {
	name, err := p.expectToken(Identifier)
	if err != nil {
		return &NoOpNode{}, err
	}
	typ, err := p.parseType()
	if err != nil {
		return &NoOpNode{}, err
	}
	return &ParameterNode{name: name.str, typ: typ}, nil
}

func (p *Parser) parseParameterList() (Node, error) {
	var paramList []Node
	for p.currentToken().kind != CloseParen {
		param, err := p.parseParameter()
		if err != nil {
			return &NoOpNode{}, err
		}
		paramList = append(paramList, param)
		switch p.currentToken().kind {
		case Comma:
			p.consumeToken()
		case CloseParen:
			break
		default:
			return &NoOpNode{}, fmt.Errorf("Parameter declaration must be followed by ) or ,")
		}
	}
	return &ParameterListNode{parameters: paramList}, nil
}

func (p *Parser) parseReturnType() (Type, error) {
	typ, err := p.parseType()
	if err != nil {
		return TypeUndetermined, err
	}
	if p.currentToken().kind != OpenCurly {
		return TypeUndetermined, fmt.Errorf("Expected `{` after return type declaration")
	}
	return typ, nil

}

func (p *Parser) parseFunction() (Node, error) {
	_, err := p.expectToken(Keyword) // fn
	if err != nil {
		return &NoOpNode{}, err
	}
	functionName, err := p.expectToken(Identifier) // function name
	if err != nil {
		return &NoOpNode{}, err
	}

	_, err = p.expectToken(OpenParen)
	if err != nil {
		return &NoOpNode{}, err
	}

	parameterList, err := p.parseParameterList()
	if err != nil {
		return &NoOpNode{}, err
	}

	_, err = p.expectToken(CloseParen)
	if err != nil {
		return &NoOpNode{}, err
	}

	returnType := NoReturnType
	if p.currentToken().kind == RightArrow {
		p.consumeToken()
		returnType, err = p.parseReturnType()
		if err != nil {
			return &NoOpNode{}, err
		}
	}

	isNew := p.createFunctionInCurrentScope(functionName.str, parameterList.(*ParameterListNode), returnType)
	if !isNew {
		return &NoOpNode{}, fmt.Errorf("Function with name %q already exists in the same scope", functionName.str)
	}

	functionBody, err := p.parseCompoundStatement(parameterList.(*ParameterListNode), returnType)
	if err != nil {
		return &NoOpNode{}, err
	}

	return &FunctionNode{name: functionName, parameters: parameterList, body: functionBody, returnType: returnType}, nil
}

func (p *Parser) parseArgumentList() ([]Node, error) {
	var arguments []Node

	_, err := p.expectToken(OpenParen)
	if err != nil {
		return arguments, err
	}

	for p.currentToken().kind != CloseParen {
		argument, err := p.parseExpr()
		if err != nil {
			return arguments, err
		}
		arguments = append(arguments, argument)
		switch p.currentToken().kind {
		case Comma:
			p.consumeToken()
		case CloseParen:
			break
		default:
			return arguments, fmt.Errorf("Expected , or ) after argument in argument list")
		}
	}
	return arguments, nil
}

func (p *Parser) parseFunctionCall() (Node, error) {
	var functionNode Token
	var err error
	switch p.currentToken().kind {
	case Identifier:
		// TODO: Check if function is defined
		functionNode, err = p.expectToken(Identifier)

	// Special case for reserved keywords that are called as functions
	case Keyword:
		switch keyword := p.currentToken().str; keyword {
		case "print":
			p.addImport("fmt")
			functionNode, err = p.expectToken(Keyword)
		default:
			return &NoOpNode{}, fmt.Errorf("UNREACHABLE: Unsupported keyword in function call: %q", keyword)
		}
	}

	if err != nil {
		return &NoOpNode{}, err
	}

	argumentList, err := p.parseArgumentList()
	// TODO: Check if right number of arguments
	if err != nil {
		return &NoOpNode{}, err
	}

	_, err = p.expectToken(CloseParen)
	if err != nil {
		return &NoOpNode{}, err
	}
	return &FunctionCallNode{name: functionNode.str, arguments: argumentList}, nil
}

func (p *Parser) parseReturn() (Node, error) {

	_, err := p.expectToken(Keyword) // return
	if err != nil {
		return &NoOpNode{}, err
	}

	expr, err := p.parseExpr()
	if err != nil {
		return &NoOpNode{}, err
	}

	return &ReturnNode{expr: expr}, nil
}

func (p *Parser) parseStatement() (Node, error) {
	switch p.currentToken().kind {

	case Identifier:
		switch p.peek(1).kind {
		case Assign:
			node, err := p.parseAssign()
			if err != nil {
				return &NoOpNode{}, err
			}
			return node, nil
		case OpenParen:
			node, err := p.parseFunctionCall()
			if err != nil {
				return &NoOpNode{}, err
			}
			return node, nil
		default:
			return  &NoOpNode{}, fmt.Errorf("TODO: Identifier follow by non-assignment in statement")
		}

	case Keyword:
		switch p.currentToken().str {
		case "fn":
			node, err := p.parseFunction()
			if err != nil {
				return &NoOpNode{}, err
			}
			return node, nil
		case "print":
			node, err := p.parseFunctionCall()
			if err != nil {
				return &NoOpNode{}, err
			}
			return node, nil
		case "return":
			node, err := p.parseReturn()
			if err != nil {
				return &NoOpNode{}, err
			}
			return node, nil
		default:
			return  &NoOpNode{}, fmt.Errorf("TODO: Parsing of keyword %q not implemented", p.currentToken().str)
		}

	case OpenCurly:
		node, err := p.parseCompoundStatement(nil, NoReturnType)
		if err != nil {
			return &NoOpNode{}, err
		}
		return node, nil
	case Integer, Float:
		node, err := p.parseExpr()
		if err != nil {
			return &NoOpNode{}, err
		}
		return node, nil
	default:
		return &NoOpNode{}, fmt.Errorf("Unknown first token in statement: %q", p.currentToken().kind)
	}
}

func (p *Parser) parseVar(checkIfDeclared bool) (Node, error) {
	token, err := p.expectToken(Identifier)
	if err != nil {
		return &NoOpNode{}, err
	}

	if checkIfDeclared {
		if isDeclared := p.validateVariable(token.str); !isDeclared {
			return &NoOpNode{}, fmt.Errorf("Use of non-declared variable: %q", token.str)
		}
	}
	return &VarNode{token: token}, nil

}

func (p *Parser) parseAssign() (Node, error) {
	left, err := p.parseVar(false)
	if err != nil {
		return &NoOpNode{}, err
	}
	isNew := p.createVariableInCurrentScope(left.(*VarNode).token.str, TypeUndetermined)
	token, err := p.expectToken(Assign)
	if err != nil {
		return &NoOpNode{}, err
	}
	right, err := p.parseExpr()
	if err != nil {
		return &NoOpNode{}, err
	}
	return &AssignNode{left: left, tok: token, right: right, declaration: isNew}, nil
}

func Parse(tokens []Token) (Node, error) {
	rootScope := newScope(nil, nil, NoReturnType)
	parser := Parser{tokens, 0, 0, rootScope, make(map[string]bool)}

	var functions []Node
	for parser.currentToken().kind != Eof {
		fn, err := parser.parseFunction()
		if err != nil {
			return &ProgramNode{}, err
		}
		functions = append(functions, fn)
	}
	return &ProgramNode{functions: functions, imports: parser.imports}, nil
}
