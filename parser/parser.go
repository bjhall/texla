package parser

import (
	"fmt"
)

type Symbol struct {
	typ        Type
	name       string
	used       bool
	fallible   bool
	category   SymbolCategory
	paramsNode *ParameterListNode
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
	fallible   bool
}

type SymbolCategory int

const (
	VariableSymbol SymbolCategory = iota
	FunctionSymbol
)

func newScope(parent *Scope, parameters []ParameterNode, returnType Type, fallible bool) *Scope {

	symbols := make(map[string]Symbol)

	// Add function parameters to the scopes list of declared symbols
	if parameters != nil {
		for _, param := range parameters {
			symbols[param.name] = Symbol{param.typ, param.name, false, false, VariableSymbol, &ParameterListNode{}}
		}
	}

	return &Scope{parent, symbols, returnType, fallible}
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

func (s *Scope) createSymbol(name string, category SymbolCategory, typ Type, paramsNode *ParameterListNode, fallible bool) bool {
	if _, exists := s.symbols[name]; exists {
		return false
	}
	s.symbols[name] = Symbol{typ, name, false, fallible, category, paramsNode}
	return true
}

func (s *Scope) closestReturningScope() *Scope {
	rs := s
	for rs.returnType == (NoReturn{}) {
		rs = rs.parent
	}
	return rs
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
	return p.currentScope.createSymbol(name, VariableSymbol, typ, &ParameterListNode{}, false)
}

func (p *Parser) createFunctionInCurrentScope(name string, paramsNode *ParameterListNode, returnType Type, fallible bool) bool {
	return p.currentScope.createSymbol(name, FunctionSymbol, returnType, paramsNode, fallible)
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

func (p *Parser) newScope(parameters []ParameterNode, returnType Type, fallible bool) {
	p.currentScope = newScope(p.currentScope, parameters, returnType, fallible)
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

	for p.currentToken().kind == Equal || p.currentToken().kind == NotEqual || p.currentToken().kind == Greater || p.currentToken().kind == GreaterEqual || p.currentToken().kind == Less || p.currentToken().kind == LessEqual {
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

		node := &NumNode{token: p.consumeToken()}

		// Chained function call
		if p.currentToken().kind == Period {
			p.consumeToken() // .
			chained, err := p.parseFunctionCall(node)
			if err != nil {
				return &NoOpNode{}, err
			}
			return chained, nil
		}

		return node, nil

	case Identifier:
		switch p.peek(1).kind {
		case OpenParen: // Function call
			functionCall, err := p.parseFunctionCall(nil)
			if err != nil {
				return &NoOpNode{}, err
			}
			return functionCall, nil
		default: // Variable
			variable, err := p.parseVar(true)
			if err != nil {
				return &NoOpNode{}, err
			}

			// Chained function call
			if p.currentToken().kind == Period {
				p.consumeToken() // .
				chained, err := p.parseFunctionCall(variable)
				if err != nil {
					return &NoOpNode{}, err
				}
				return chained, nil
			}

			return variable, nil
		}

	case Keyword:
		switch keyword := p.currentToken().str; keyword {
		case "true":
			return &BoolNode{token: p.consumeToken()}, nil
		case "false":
			return &BoolNode{token: p.consumeToken()}, nil
		default:
			return &NoOpNode{}, fmt.Errorf("Unexpected keyword: %q", keyword)
		}

	case StringLiteral:

		node := &StringLiteralNode{token: p.consumeToken()}

		// Chained function call
		if p.currentToken().kind == Period {
			p.consumeToken() // .
			chained, err := p.parseFunctionCall(node)
			if err != nil {
				return &NoOpNode{}, err
			}
			return chained, nil
		}

		return node, nil

	case Not:
		op := p.consumeToken()
		expr, err := p.parseExpr()
		if err != nil {
			return &NoOpNode{}, err
		}
		return &UnaryOpNode{op: op, expr: expr}, nil

	case OpenBracket:
		slice, err := p.parseSliceLiteral()
		if err != nil {
			return &NoOpNode{}, err
		}
		return slice, nil

	default:
		return &NoOpNode{}, fmt.Errorf("Unexpected token: %q", p.currentToken().str)
	}
}

func (p *Parser) parseSliceLiteral() (Node, error) {
	_, err := p.expectToken(OpenBracket)
	if err != nil {
		return &NoOpNode{}, err
	}

	// TODO: Parsing a comma separated list of expression should probably be extracted to a function
	var elements []Node
	for {
		expr, err := p.parseExpr()
		if err != nil {
			return &NoOpNode{}, err
		}
		elements = append(elements, expr)
		if p.currentToken().kind == CloseBracket || p.currentToken().kind == Eof {
			break
		}
		_, err = p.expectToken(Comma)
		if err != nil {
			return &NoOpNode{}, err
		}

		// Allow slice literal to end with comma
		if p.currentToken().kind == CloseBracket {
			break
		}
	}

	_, err = p.expectToken(CloseBracket)
	if err != nil {
		return &NoOpNode{}, err
	}

	return &SliceLiteralNode{elements: elements}, nil
}

func (p *Parser) parseCompoundStatement(parameters []ParameterNode, returnType Type, fallible bool) (Node, error) {

	_, err := p.expectToken(OpenCurly)
	if err != nil {
		return &NoOpNode{}, err
	}

	// Create new scope
	p.newScope(parameters, returnType, fallible)
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
		return TypeUndetermined{}, err
	}
	switch typeToken.str {
	case "int":
		return TypeInt{}, nil
	case "float":
		return TypeFloat{}, nil
	case "str":
		return TypeString{}, nil
	case "bool":
		return TypeBool{}, nil
	default:
		return TypeUndetermined{}, fmt.Errorf("Unknown type: %q", typeToken.str)
	}
}

func literalTokenType(tok Token) (Type, error) {
	switch tok.kind {
	case Integer:
		return TypeInt{}, nil
	case Float:
		return TypeFloat{}, nil
	case StringLiteral:
		return TypeString{}, nil
	case Keyword:
		if tok.str == "true" || tok.str == "false" {
			return TypeBool{}, nil
		}
		return TypeUndetermined{}, fmt.Errorf("Token is not literal %q", tok.kind)
	default:
		fmt.Println(tok)
		return TypeUndetermined{}, fmt.Errorf("Token is not literal %q", tok.kind)
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
	if p.currentToken().kind == Assign {
		p.consumeToken()
		defaultToken := p.consumeToken()
		literalType, err := literalTokenType(defaultToken)
		if err != nil {
			return &NoOpNode{}, err
		}
		if typ == literalType {
			return &ParameterNode{name: name.str, typ: typ, hasDefault: true, defaultValue: defaultToken.str}, nil
		}
		return &NoOpNode{}, fmt.Errorf("Default argument has wrong type")
	}
	return &ParameterNode{name: name.str, typ: typ}, nil
}

func (p *Parser) parseParameterList() (Node, error) {
	var paramList []ParameterNode
	for p.currentToken().kind != CloseParen {
		param, err := p.parseParameter()
		if err != nil {
			return &NoOpNode{}, err
		}
		paramList = append(paramList, *param.(*ParameterNode))
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
		return TypeUndetermined{}, err
	}
	if p.currentToken().kind != OpenCurly {
		return TypeUndetermined{}, fmt.Errorf("Expected `{` after return type declaration")
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

	fallible := false
	if p.currentToken().kind == QuestionMark {
		fallible = true
		p.consumeToken() // ?
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

	var returnType Type
	if p.currentToken().kind == RightArrow {
		p.consumeToken()
		returnType, err = p.parseReturnType()
		if err != nil {
			return &NoOpNode{}, err
		}
	} else {
		returnType = TypeVoid{}
	}

	isNew := p.createFunctionInCurrentScope(functionName.str, parameterList.(*ParameterListNode), returnType, fallible)
	if !isNew {
		return &NoOpNode{}, fmt.Errorf("Function with name %q already exists in the same scope", functionName.str)
	}

	functionBody, err := p.parseCompoundStatement(parameterList.(*ParameterListNode).parameters, returnType, fallible)
	if err != nil {
		return &NoOpNode{}, err
	}

	return &FunctionNode{name: functionName, parameters: parameterList, body: functionBody, returnType: returnType, fallible: fallible}, nil
}

func (p *Parser) parseArgumentList(self Node) ([]Node, error) {
	var arguments []Node

	_, err := p.expectToken(OpenParen)
	if err != nil {
		return arguments, err
	}

	namedArgumentSeen := false
	orderedArgumentCount := 0

	// If function chaining was used, add "self" argument to beginning of argument list
	if self != nil {
		arguments = append(arguments, &ArgumentNode{expr: self, named: false, order: 0})
		orderedArgumentCount++
	}

	for p.currentToken().kind != CloseParen {

		paramName := ""

		// Check if named argument
		if p.currentToken().kind == Identifier && p.peek(1).kind == Assign {
			paramName = p.consumeToken().str
			p.consumeToken() // =
			namedArgumentSeen = true
		} else if namedArgumentSeen {
			return arguments, fmt.Errorf("Named argument cannot be followed by unnamed argument in list of arguments")
		}

		// Parse the actual argument
		argumentExpr, err := p.parseExpr()
		if err != nil {
			return arguments, err
		}

		if paramName != "" { // Named argument
			arguments = append(arguments, &ArgumentNode{expr: argumentExpr, named: true, paramName: paramName})
		} else { // Ordered argument
			arguments = append(arguments, &ArgumentNode{expr: argumentExpr, named: false, order: orderedArgumentCount})
			orderedArgumentCount++
		}
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

func (p *Parser) parseFunctionCall(self Node) (Node, error) {
	var functionToken Token
	var err error
	switch p.currentToken().kind {
	case Identifier:
		functionToken, err = p.expectToken(Identifier)

	// Special case for reserved keywords that are called as functions
	case Keyword:
		switch keyword := p.currentToken().str; keyword {
		case "print":
			p.addImport("fmt")
			functionToken, err = p.expectToken(Keyword)
		default:
			return &NoOpNode{}, fmt.Errorf("UNREACHABLE: Unsupported keyword in function call: %q", keyword)
		}
	}

	if err != nil {
		return &NoOpNode{}, err
	}

	argumentList, err := p.parseArgumentList(self)

	if err != nil {
		return &NoOpNode{}, err
	}

	_, err = p.expectToken(CloseParen)
	if err != nil {
		return &NoOpNode{}, err
	}

	errorHandled := false
	if p.currentToken().kind == QuestionMark {
		errorHandled = true
		p.consumeToken() // ?
		if p.currentToken().kind == OpenCurly {
			errVariable := &ParameterNode{name: "err", typ: TypeString{}}
			errBody, err := p.parseCompoundStatement([]ParameterNode{*errVariable}, NoReturn{}, false)
			if err != nil {
				return &NoOpNode{}, err
			}
			return &FunctionCallNode{name: functionToken.str, arguments: argumentList, errorHandled: true, errorBody: errBody}, nil
		}
	}

	// Special case for generator build-ins such as read()
	if isBuiltin(functionToken.str) && p.currentToken().kind == RightArrow {

		p.consumeToken() // ->

		_, isGenerator := builtins[functionToken.str].returnType.(*TypeGenerator)
		if isGenerator {
			return &NoOpNode{}, fmt.Errorf("Cannot put -> after non-generator function %q", functionToken.str)
		}

		variable, err := p.parseVar(false)
		if err != nil {
			return &NoOpNode{}, err
		}
		controlVariable := &ParameterNode{name: variable.(*VarNode).token.str, typ: TypeUndetermined{}}

		body, err := p.parseCompoundStatement([]ParameterNode{*controlVariable}, NoReturn{}, false)
		//_ = p.createVariableInCurrentScope(lhsName, TypeUndetermined{})
		if err != nil {
			return &NoOpNode{}, err
		}

		variableNode, _ := variable.(*VarNode)
		return &FunctionCallNode{name: functionToken.str, arguments: argumentList, isBuiltin: isBuiltin(functionToken.str), errorHandled: errorHandled, generatorVar: *variableNode, generatorBody: body}, nil
	}

	functionCall := &FunctionCallNode{name: functionToken.str, arguments: argumentList, isBuiltin: isBuiltin(functionToken.str), errorHandled: errorHandled}

	// Chained function call
	if p.currentToken().kind == Period {
		p.consumeToken() // .
		chained, err := p.parseFunctionCall(functionCall)
		if err != nil {
			return &NoOpNode{}, err
		}
		return chained, nil
	}

	return functionCall, nil
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

func (p *Parser) parseFail() (Node, error) {

	_, err := p.expectToken(Keyword) // fail
	if err != nil {
		return &NoOpNode{}, err
	}

	expr, err := p.parseExpr()
	if err != nil {
		return &NoOpNode{}, err
	}

	return &FailNode{expr: expr}, nil
}

func (p *Parser) parseIfStatement() (Node, error) {
	_, err := p.expectToken(Keyword) // if
	if err != nil {
		return &NoOpNode{}, err
	}

	comp, err := p.parseComparison()
	if err != nil {
		return &NoOpNode{}, err
	}

	body, err := p.parseCompoundStatement(nil, NoReturn{}, false)
	if err != nil {
		return &NoOpNode{}, err
	}

	if p.currentToken().kind == Keyword && p.currentToken().str == "else" {
		p.consumeToken() // else

		var elseBody Node
		// else if
		if p.currentToken().kind == Keyword && p.currentToken().str == "if" {
			elseBody, err = p.parseIfStatement()
		} else { // just else
			elseBody, err = p.parseCompoundStatement(nil, NoReturn{}, false)
		}
		if err != nil {
			return &NoOpNode{}, err
		}
		return &IfNode{comp: comp, body: body, elseBody: elseBody}, nil
	}

	return &IfNode{comp: comp, body: body, elseBody: &NoOpNode{}}, nil
}

func (p *Parser) parseIterator() (Node, error) {
	firstExpr, err := p.parseExpr()
	if err != nil {
		return &NoOpNode{}, err
	}
	switch p.currentToken().kind {
	case Range:
		rangeToken := p.consumeToken() // ..
		rangeEnd, err := p.parseExpr()
		if err != nil {
			return &NoOpNode{}, err
		}
		return &RangeNode{token: rangeToken, from: firstExpr, to: rangeEnd, step: 1}, nil

	default:
		switch firstExpr.Type() {
		case VarNodeType, SliceLiteralNodeType:
			return firstExpr, nil
		default:
			panic("Non-supported iterator...")
		}
	}
}

func (p *Parser) parseForLoop() (Node, error) {
	_, err := p.expectToken(Keyword) // for
	if err != nil {
		return &NoOpNode{}, err
	}

	iterator, err := p.parseIterator()
	if err != nil {
		return &NoOpNode{}, err
	}
	_, isRange := iterator.(*RangeNode)

	_, err = p.expectToken(RightArrow) // ->
	if err != nil {
		return &NoOpNode{}, err
	}

	variable, err := p.parseVar(false)
	if err != nil {
		return &NoOpNode{}, err
	}
	controlVariable := &ParameterNode{name: variable.(*VarNode).token.str, typ: TypeUndetermined{}}
	forParams := []ParameterNode{*controlVariable}

	// Parse optional index variable
	var idxVariable Node
	hasIdx := false
	if p.currentToken().kind == Comma && !isRange{
		p.consumeToken() // ,
		idxVariable, err = p.parseVar(false)
		if err != nil {
			return &NoOpNode{}, err
		}
		hasIdx = true
		idxVariableParam := &ParameterNode{name: idxVariable.(*VarNode).token.str, typ: TypeInt{}}
		forParams = append(forParams, *idxVariableParam)
	}

	body, err := p.parseCompoundStatement(forParams, NoReturn{}, false)
	if err != nil {
		return &NoOpNode{}, err
	}

	variableNode, _ := variable.(*VarNode)
	if hasIdx {
		idxVariableNode, _ := idxVariable.(*VarNode)
		return &ForeachNode{iterator: iterator, variable: *variableNode, idxVariable: *idxVariableNode, body: body, hasIdx: hasIdx}, nil
	} else {
		return &ForeachNode{iterator: iterator, variable: *variableNode, idxVariable: VarNode{}, body: body, hasIdx: hasIdx}, nil
	}
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
			node, err := p.parseFunctionCall(nil)
			if err != nil {
				return &NoOpNode{}, err
			}
			return node, nil
		case Period: // FIXME: This was added to allow chained function calls as statements, eg `a.append(1)`. Is it correct?
			node, err := p.parseExpr()
			if err != nil {
				return &NoOpNode{}, err
			}
			return node, nil
		default:
			return &NoOpNode{}, fmt.Errorf("Syntax error (ADD PROPER ERROR MESSAGE HERE")
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
			node, err := p.parseFunctionCall(nil)
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
		case "if":
			node, err := p.parseIfStatement()
			if err != nil {
				return &NoOpNode{}, err
			}
			return node, nil
		case "for":
			node, err := p.parseForLoop()
			if err != nil {
				return &NoOpNode{}, err
			}
			return node, nil
		case "true", "false":
			node, err := p.parsePrimary()
			if err != nil {
				return &NoOpNode{}, err
			}
			return node, nil

		case "fail":
			node, err := p.parseFail()
			if err != nil {
				return &NoOpNode{}, err
			}
			return node, nil
		default:
			return &NoOpNode{}, fmt.Errorf("TODO: Parsing of keyword %q not implemented", p.currentToken().str)
		}

	case OpenCurly:
		node, err := p.parseCompoundStatement(nil, NoReturn{}, false)
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
			return &NoOpNode{}, fmt.Errorf("Use of non-declared variable: %q %d:%d", token.str, token.line, token.column)
		}
	}

	// Indexing (eg. a[10])
	if p.currentToken().kind == OpenBracket {
		p.consumeToken() // [
		indexNode, err := p.parseExpr()
		_, err = p.expectToken(CloseBracket)
		if err != nil {
			return &NoOpNode{}, err
		}
		return &IndexedVarNode{token: token, index: indexNode}, nil
	}
	return &VarNode{token: token}, nil

}

func (p *Parser) parseAssign() (Node, error) {
	left, err := p.parseVar(false)
	if err != nil {
		return &NoOpNode{}, err
	}
	lhsName := left.(*VarNode).token.str
	_, exists := p.currentScope.lookupSymbol(lhsName)
	if !exists {
		_ = p.createVariableInCurrentScope(lhsName, TypeUndetermined{})
	}
	token, err := p.expectToken(Assign)
	if err != nil {
		return &NoOpNode{}, err
	}
	right, err := p.parseExpr()
	if err != nil {
		return &NoOpNode{}, err
	}
	if p.currentToken().kind == Range {
		// TODO: Factor this stuff out to a parseRange that takes the "from" node as an argument
		rangeToken := p.consumeToken() // ..
		end, err := p.parseExpr()
		if err != nil {
			return &NoOpNode{}, err
		}
		return &AssignNode{left: left, tok: token, right: &RangeNode{token: rangeToken, from: right, to: end, step: 1}, declaration: !exists}, nil
	}
	return &AssignNode{left: left, tok: token, right: right, declaration: !exists}, nil
}

func Parse(tokens []Token) (Node, error) {
	rootScope := newScope(nil, nil, NoReturn{}, false)
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
