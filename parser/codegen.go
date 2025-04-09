package parser

import (
	"fmt"
	"strings"
)

type Generator struct {
	indentLevel int
	scope       *Scope
	errors      []string
	imports     map[string]bool
	preludes    map[string]bool
}

func (g *Generator) addPreludeFunction(name string) {
	for _, imp := range preludeImports(name) {
		g.addImport(imp)
	}
	g.preludes[name] = true
}

func (g *Generator) addImport(name string) {
	g.imports[name] = true
}


func (g *Generator) error(errorStr string) {
	g.errors = append(g.errors, fmt.Sprintf("CODEGEN ERROR: %s", errorStr))
}


func (g *Generator) indent(str string) string {
	var indentation string
	for _ = range g.indentLevel {
		indentation += "    "
	}
	return indentation + str
}


func (g *Generator) coerce(content string, from Type, to Type, mode CoercionMode) string {
	if from == to || to == NoCoercion {
		return content
	}

	switch from {
	case TypeInt:
		switch to {
		case TypeFloat:
			return fmt.Sprintf("float64(%s)", content)
		case TypeString:
			g.addPreludeFunction("intToString")
			return fmt.Sprintf("strconv.Itoa(%s)", content)
		case TypeBool:
			return fmt.Sprintf("%s != 0", content)
		default:
			panic("Unimplemented coercion")
		}
	case TypeFloat:
		switch to {
		case TypeInt:
			if mode == CoercionModeNumLiteral {
				g.addImport("math")
				return fmt.Sprintf("int(math.Floor(%s))", content)
			} else {
				return fmt.Sprintf("int(%s)", content)
			}
		case TypeString:
			return fmt.Sprintf("strconv.FormatFloat(%s, 'f', -1, 64)", content)
		case TypeBool:
			return fmt.Sprintf("%s != 0", content)
		default:
			panic("Unimplemented coercion")
		}
	case TypeString:
		switch to {
		case TypeInt:
			g.addPreludeFunction("stringToInt")
			return fmt.Sprintf("stringToInt(%s)", content)
		case TypeFloat:
			g.addPreludeFunction("stringToFloat")
			return fmt.Sprintf("stringToFloat(%s)", content)
		case TypeBool:
			return fmt.Sprintf("len(%s) > 0", content)
		default:
			panic("Unimplemented coercion")
		}
	default:
		panic("Unimplemented coercion")
	}
}



func (g *Generator) codegenNum(node *NumNode, coercion Type) string {
	return g.coerce(node.token.str, node.NumType(), coercion, CoercionModeNumLiteral)
}

func (g *Generator) codegenBool(node *BoolNode) string {
	// TODO: Should booleans ever coerce?
	return node.token.str
}

func (g *Generator) codegenStringLiteral(node *StringLiteralNode, coercion Type) string {
	return g.coerce("\""+node.token.str+"\"", TypeString, coercion, CoercionModeDefault)
}


func (g *Generator) codegenVar(node *VarNode, coercion Type) string {
	varName := node.token.str
	if coercion == NoCoercion {
		return varName
	}

	symbol, _ := g.scope.lookupSymbol(varName)
	if symbol.category != VariableSymbol {
		panic("Should be variable...") // TODO: ASSERT
	}

	return g.coerce(varName, symbol.typ, coercion, CoercionModeDefault)
}


func (g *Generator) codegenUnaryOp(node *UnaryOpNode) string {
	switch node.op.kind {
	case Not:
		return fmt.Sprintf("%s(%s)", node.op.str, g.codegenExpr(node.expr, TypeBool))
	default:
		panic("Codegen for unary op not implemeneted")
	}
}

func (g *Generator) codegenBinOp(node *BinOpNode, coercion Type) string {
	left := g.codegenWithParens(node.left, node, coercion)
	right := g.codegenWithParens(node.right, node, coercion)
	return fmt.Sprintf("%s %s %s", left, node.op.str, right)
}


func (g *Generator) codegenWithParens(node Node, parent Node, coercion Type) string {
    result := g.codegenExpr(node, coercion)

    // Add parentheses if:
    // 1. The child is a binary operation
    // 2. The child's precedence is lower than the parent's
    needsParens := false

    if node.Type() == BinOpNodeType && parent.Type() == BinOpNodeType {
		childOp := node.(*BinOpNode)
		parentOp := parent.(*BinOpNode)

		// Add parens if child precedence is lower
		if childOp.Precedence() < parentOp.Precedence() {
			needsParens = true
		}

		// Special case: For right-associative operators or
		// when precedences are equal but on the right side
		if childOp.Precedence() == parentOp.Precedence() {
			if parent.(*BinOpNode).right == node {
				needsParens = true
			}
		}
    }

    if needsParens {
        return "(" + result + ")"
    }
    return result
}

func (g *Generator) codegenAssign(node *AssignNode) string {
	opStr := "="
	if node.declaration {
		opStr = ":="
	}
	lhs := g.codegenVar(node.left.(*VarNode), NoCoercion)
	lhsSymbol := g.scope.symbols[lhs]

	// If variable is not used, add `_ = varname` after assignment
	unusedStr := ""
	if !lhsSymbol.used {
		unusedStr = fmt.Sprintf("\n%s = %s\n", g.indent("_"), lhs)
	}

	return fmt.Sprintf(
		"%s %s %s%s",
		lhs,
		opStr,
		g.codegenExpr(node.right, lhsSymbol.typ),
		unusedStr,
	)
}

func (g *Generator) codegenCompoundStatement(node *CompoundStatementNode) string {
	prevScope := g.scope
	g.scope = node.scope
	var statements []string
	g.indentLevel++
	for _, child := range node.children {
		statements = append(statements, g.indent(g.codegenStatement(child)))
	}
	statementsString := strings.Join(statements, "\n")
	g.indentLevel--
	g.scope = prevScope

	return fmt.Sprintf("{\n%s\n%s", statementsString, g.indent("}"))
}

func (g *Generator) codegenType(typ Type) string {
	switch typ {
	case TypeInt:
		return "int"
	case TypeFloat:
		return "float64"
	case TypeString:
		return "string"
	default:
		panic("UNIMPLEMENTED TYPE") // FIXME
	}
}

func (g *Generator) codegenReturnType(typ Type) string {
	if typ == NoReturnType {
		return ""
	}
	return g.codegenType(typ)
}

func (g *Generator) codegenParameter(node *ParameterNode) string {
	return fmt.Sprintf("%s %s", node.name, g.codegenType(node.typ))
}

func (g *Generator) codegenParameterList(node *ParameterListNode) string {
	var parameters []string
	for _, param := range node.parameters {
		parameters = append(parameters, g.codegenParameter(param.(*ParameterNode)))
	}
	return strings.Join(parameters, ", ")
}

func (g *Generator) codegenFunction(node *FunctionNode) string {
	return fmt.Sprintf("func %s(%s) %s %s",
		node.name.str,
		g.codegenParameterList(node.parameters.(*ParameterListNode)),
		g.codegenReturnType(node.returnType),
		g.codegenCompoundStatement(node.body.(*CompoundStatementNode)),
	)
}

func (g *Generator) codegenFunctionCall(node *FunctionCallNode, coercion Type) string {

	var functionName string
	var symbol Symbol
	var argumentStrings []string

	// Some built-ins, such as `print` are essentially just renames of other functions, that's handled here
	switch node.name {
	case "print":
		functionName = "fmt.Println"
	default:
		functionName = node.name
		symbol, _ = g.scope.lookupSymbol(functionName)
	}

	if node.name == "print" {
		for _, argument := range node.arguments {
			argumentStrings = append(argumentStrings, g.codegenExpr(argument, NoCoercion))
		}
	} else {
		for i, argument := range node.arguments {
			parameterType := symbol.parameterTypes[i]
			//argumentStrings = append(argumentStrings, convertedArguments[i])
			argumentStrings = append(argumentStrings, g.codegenExpr(argument, parameterType))
		}
	}
	functionCall := fmt.Sprintf("%s(%s)", functionName, strings.Join(argumentStrings, ", "))
	return g.coerce(functionCall, symbol.typ, coercion, CoercionModeDefault)
}

func (g *Generator) codegenReturn(node *ReturnNode) string {
	expectedReturnType := g.scope.returnType
	return fmt.Sprintf("return %s", g.codegenExpr(node.expr, expectedReturnType))
}

func (g *Generator) codegenIf(node *IfNode) string {
	coerceType := node.compType
	if node.comp.Type() == VarNodeType {
		coerceType = TypeBool
	}

	return fmt.Sprintf(
		"if %s %s",
		g.codegenExpr(node.comp, coerceType),
		g.codegenCompoundStatement(node.body.(*CompoundStatementNode)),
	)
}

func (g *Generator) codegenStatement(node Node) string {
	switch node.Type() {
	case AssignNodeType:
		return g.codegenAssign(node.(*AssignNode))
	case CompoundStatementNodeType:
		return g.codegenCompoundStatement(node.(*CompoundStatementNode))
	case FunctionNodeType:
		return g.codegenFunction(node.(*FunctionNode))
	case FunctionCallNodeType:
		return g.codegenFunctionCall(node.(*FunctionCallNode), NoCoercion)
	case ReturnNodeType:
		return g.codegenReturn(node.(*ReturnNode))
	case IfNodeType:
		return g.codegenIf(node.(*IfNode))
	default:
		fmt.Println("CODEGEN TODO: Unknown node in statement", node.Type())
		panic("")
	}
}

func (g *Generator) codegenExpr(node Node, coercion Type) string {
	switch node.Type() {
	case UnaryOpNodeType:
		return g.codegenUnaryOp(node.(*UnaryOpNode))
	case BinOpNodeType:
		return g.codegenBinOp(node.(*BinOpNode), coercion)
	case NumNodeType:
		return g.codegenNum(node.(*NumNode), coercion)
	case BoolNodeType:
		return g.codegenBool(node.(*BoolNode))
	case StringLiteralNodeType:
		return g.codegenStringLiteral(node.(*StringLiteralNode), coercion)
	case VarNodeType:
		return g.codegenVar(node.(*VarNode), coercion)
	case FunctionCallNodeType:
		return g.codegenFunctionCall(node.(*FunctionCallNode), coercion)
	default:
		fmt.Println("CODEGEN TODO: Unknown node in expression", node.Type())
		panic("")
	}
}

func (g *Generator) codegenProgram(node Node) string {
	var functionStrs []string
	for _, function := range node.(*ProgramNode).functions {
		functionStrs = append(functionStrs, g.codegenFunction(function.(*FunctionNode)))
	}
	var importStrs []string
	for imp, _ := range g.imports {
		node.(*ProgramNode).addImport(imp)
	}
	for imp, _ := range node.(*ProgramNode).imports {
		importStrs = append(importStrs, "import \""+imp+"\"")
	}

	prelude := ""

	for preludeName, _ := range g.preludes {
		prelude += preludeCode(preludeName)
	}

	return fmt.Sprintf("package main\n\n%s\n\n%s\n\n%s", strings.Join(importStrs, "\n"), prelude, strings.Join(functionStrs, "\n\n"))


}

func GenerateCode(root Node) (string, error) {
	generator := Generator{0, nil, []string{}, make(map[string]bool), make(map[string]bool)}
	code := generator.codegenProgram(root)

	if len(generator.errors) > 0 {
		return "", fmt.Errorf(strings.Join(generator.errors, "\n"))
	}

	return code, nil
}
