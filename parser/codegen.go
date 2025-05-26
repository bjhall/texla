package parser

import (
	"fmt"
	"strings"
)

type Generator struct {
	indentLevel         int
	scope               *Scope
	errors              []string
	imports             map[string]bool
	preludes            map[string]bool
	initStatements      []string
	finalStatements     []string
	preStatements       []string
	postStatements      []string
	replacementCount    int
	ignorePreStatements bool
	tmpVarCount         int
}


func (g *Generator) nilValue(typ Type) string {
	switch typ.(type) {
	case TypeInt, TypeFloat:
		return "0"
	case TypeString:
		return "\"\""
	case TypeSlice:
		return typ.String()+"{}"
	default:
		panic("TODO: Unimplemented nil value for type in fail")
	}
}

func (g *Generator) getReplacementVarName(fnName string) string {
	g.replacementCount++
	return fmt.Sprintf("___%s_result_%d", fnName, g.replacementCount)
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

func (g *Generator) addInitStatement(node string) {
	g.initStatements = append(g.initStatements, node)
}

func (g *Generator) addFinalStatement(node string) {
	g.finalStatements = append(g.finalStatements, node)
}

func (g *Generator) addPreStatement(node string) {
	g.preStatements = append(g.preStatements, node)
}

func (g *Generator) addPostStatement(node string) {
	g.postStatements = append(g.postStatements, node)
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
	if from == to || to.String() == "NoCoercion" {
		return content
	}

	switch from.(type) {
	case TypeInt:
		switch to.(type) {
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
		switch to.(type) {
		case TypeInt:
			if mode == CoercionModeNumLiteral {
				g.addImport("math")
				return fmt.Sprintf("int(math.Floor(%s))", content)
			} else {
				return fmt.Sprintf("int(%s)", content)
			}
		case TypeString:
			g.addImport("strconv")
			return fmt.Sprintf("strconv.FormatFloat(%s, 'f', -1, 64)", content)
		case TypeBool:
			return fmt.Sprintf("%s != 0", content)
		default:
			panic("Unimplemented coercion")
		}
	case TypeString:
		switch to.(type) {
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
	case TypeSlice:
		switch to.(type) {
		case TypeBool:
			return fmt.Sprintf("len(%s) > 0", content)
		case TypeInt:
			return fmt.Sprintf("len(%s) > 0", content) // This was added to for example `if str.find("something") {` work. Does it cause any unwanted side effects?
		default:
			panic("Unimplemented coercion for slice")
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
	return g.coerce("\""+node.token.str+"\"", TypeString{}, coercion, CoercionModeDefault)
}

func literalToStr(value string, typ Type) string {
	switch typ.(type) {
	case TypeInt, TypeFloat, TypeBool:
		return value
	case TypeString:
		return fmt.Sprintf("\"%s\"", value)
	default:
		panic("Type not supported")
	}
}

func (g *Generator) codegenVar(node *VarNode, coercion Type) string {
	varName := node.token.str

	if coercion.String() == "NoCoercion" {
		return varName
	}

	symbol, _ := g.scope.lookupSymbol(varName)
	if symbol.category != VariableSymbol {
		panic("Should be variable...") // TODO: ASSERT
	}

	return g.coerce(varName, symbol.typ, coercion, CoercionModeDefault)
}

func (g *Generator) codegenIndexedVar(node *IndexedVarNode, coercion Type) string {
	varName := node.token.str

	indexedVar := fmt.Sprintf("%s[%s]", varName, g.codegenExpr(node.index, TypeInt{}))

	symbol, _ := g.scope.lookupSymbol(varName)
	if symbol.category != VariableSymbol {
		panic("Should be variable...") // TODO: ASSERT
	}

	switch t := symbol.typ.(type) {
	case TypeSlice:
		return g.coerce(indexedVar, t.ElementType, coercion, CoercionModeDefault)
	case TypeString:
		return fmt.Sprintf("string(%s)", g.coerce(indexedVar, TypeString{}, coercion, CoercionModeDefault))
	default:
		panic("Non-indxable type")
	}
}

func (g *Generator) codegenSliceLiteral(node *SliceLiteralNode, coercion Type) string {

	elements := []string{}
	for _, elem := range node.elements {
		elements = append(elements, g.codegenExpr(elem, node.elementType))
	}

	return fmt.Sprintf("[]%s{%s}", g.codegenType(node.elementType), strings.Join(elements, ","))
}

func (g *Generator) codegenUnaryOp(node *UnaryOpNode) string {
	switch node.op.kind {
	case Not:
		return fmt.Sprintf("%s(%s)", node.op.str, g.codegenExpr(node.expr, TypeBool{}))
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

func (g *Generator) codegenAssignExpr(node *AssignNode, coercion Type) string {
	if !node.expression {
		panic("Assignment not expression")
	}
	preAssignment := g.codegenAssign(node)
	g.addPreStatement(preAssignment)
	return g.codegenVar(node.left.(*VarNode), coercion)
}

func (g *Generator) codegenAssign(node *AssignNode) string {
	opStr := "="
	if node.declaration {
		opStr = ":="
	}
	lhs := g.codegenVar(node.left.(*VarNode), NoCoercion{})
	lhsSymbol, found := g.scope.lookupSymbol(node.left.(*VarNode).token.str)
	if !found {
		panic("Codegen of non-defined symbol in assignment")
	}

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
	g.indentLevel++

	var statements []string
	for _, initStatement := range g.initStatements {
		statements = append(statements, g.indent(initStatement))
	}
	g.initStatements = nil

	// Collect final statements before codegening the body, will be added to the end later.
	// This is necessary to avoid it adding the final statements to any child scopes.
	var finals []string
	for _, finalStatement := range g.finalStatements {
		finals = append(finals, g.indent(finalStatement))
	}
	g.finalStatements = nil

	for _, child := range node.children {
		if !g.ignorePreStatements {
			g.preStatements = nil
		}

		statement := g.codegenStatement(child)

		// Handle any pre-statements generated by the codegen of the main statement
		if !g.ignorePreStatements {
			for _, preStatement := range g.preStatements {
				statements = append(statements, g.indent(preStatement))
			}
		}

		// Insert the main statement
		statements = append(statements, g.indent(statement))
	}

	// Add post statements for scopes that return values
	if g.scope.returnType != (NoReturn{}) {
		for _, postStatement := range g.postStatements {
			statements = append(statements, g.indent(postStatement))
		}
		g.postStatements = nil
	}

	// Add final statements to the end of the scope
	statements = append(statements, finals...)

	statementsString := strings.Join(statements, "\n")
	g.indentLevel--
	g.scope = prevScope

	return fmt.Sprintf("{\n%s\n%s", statementsString, g.indent("}"))
}

func (g *Generator) codegenType(typ Type) string {
	switch typ.(type) {
	case TypeInt:
		return "int"
	case TypeFloat:
		return "float64"
	case TypeString:
		return "string"
	case TypeBool:
		return "bool"
	case TypeVoid:
		return ""
	default:
		panic("UNIMPLEMENTED TYPE") // FIXME
	}
}

func (g *Generator) codegenReturnType(typ Type) string {
	if typ == nil || typ == (NoReturn{}) {
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
		parameters = append(parameters, g.codegenParameter(&param))
	}
	return strings.Join(parameters, ", ")
}

func (g *Generator) codegenFunction(node *FunctionNode) string {
	returns := g.codegenReturnType(node.returnType)

	if node.fallible {
		if node.returnType == (TypeVoid{}) {
			returns = fmt.Sprintf("error")
		} else {
			returns = fmt.Sprintf("(%s, error)", returns)
		}

	}
	paramStr := g.codegenParameterList(node.parameters.(*ParameterListNode))

	if node.fallible && node.returnType == (TypeVoid{}) {
		g.addPostStatement("return nil")
	}
	bodyStr := g.codegenCompoundStatement(node.body.(*CompoundStatementNode))
	return fmt.Sprintf("func %s(%s) %s %s",	node.name.str,paramStr, returns,bodyStr)
}

func (g *Generator) codegenFunctionCall(node *FunctionCallNode, coercion Type) string {

	// Separate codegen function for builtin calls
	if node.isBuiltin {
		return g.codegenBuiltinCall(node, coercion)
	}

	// FIXME: Make print a builtin!
	if node.name == "print" {
		var argumentStrings []string
		for _, argument := range node.arguments {
			argumentStrings = append(argumentStrings, g.codegenExpr(argument.(*ArgumentNode).expr, NoCoercion{}))
		}
		return fmt.Sprintf("fmt.Println(%s)", strings.Join(argumentStrings, ", "))

	} else {
		symbol, _ := g.scope.lookupSymbol(node.name)

		// Codegen all arguements
		var argumentStrings []string
		for _, param := range symbol.paramsNode.parameters {
			argumentStrings = append(argumentStrings, g.codegenExpr(node.resolvedArgs[param.name].expr, param.typ))
		}

		// Codegen the final call
		functionCall := fmt.Sprintf("%s(%s)", node.name, strings.Join(argumentStrings, ", "))

		// For call to non-fallible function, just return the call
		if !symbol.fallible {
			return g.coerce(functionCall, symbol.typ, coercion, CoercionModeDefault)
		}


		// For calls to fallible function, things becaome a bit more complicated...
		returnScope := g.scope.closestReturningScope()
		lhsVars := []string{"err"}
		onErrReturnVars := []string{"err"}
		replacementCode := ""
		if symbol.typ != (TypeVoid{}) {
			replacementVar := g.getReplacementVarName(node.name)
			lhsVars = []string{replacementVar, "err"}
			replacementCode = g.coerce(replacementVar, symbol.typ, coercion, CoercionModeDefault)
			if returnScope.returnType != (TypeVoid{}) {
				onErrReturnVars = []string{g.nilValue(returnScope.returnType), "err"}
			}
		}

		// Generate error-catching function call pre-statement
		g.addPreStatement(fmt.Sprintf("%s := %s", strings.Join(lhsVars, ", "), functionCall))

		// Generate error handling prestatement
		if node.errorBody != nil {
			g.ignorePreStatements = true
			g.addPreStatement(fmt.Sprintf("if err != nil %s", g.codegenCompoundStatement(node.errorBody.(*CompoundStatementNode))))
			g.ignorePreStatements = false
		} else if returnScope.fallible {
			g.addPreStatement(fmt.Sprintf("if err != nil { return %s }", strings.Join(onErrReturnVars, ", ")))
		} else {
			g.addPreludeFunction("handleNonPropagatableError")
			g.addPreStatement("___handleNonPropagatableError(err)")
		}

		return replacementCode
	}
}

func (g *Generator) codegenBuiltinCall(node *FunctionCallNode, coercion Type) string {
	builtin := builtins[node.name]

	callStr := ""
	switch builtin.name {

	case "len":
		callStr = fmt.Sprintf("len(%s)", g.codegenExpr(node.resolvedArgs["var"].expr, NoCoercion{}))

	case "append":
		destArg := node.resolvedArgs["dest"]
		dest := g.codegenVar(destArg.expr.(*VarNode), NoCoercion{})
		switch destArg.typ.(type) {
		case TypeSlice:
			return fmt.Sprintf("%s = append(%s, %s)",
				dest, dest,
				g.codegenExpr(
					node.resolvedArgs["var"].expr,
					destArg.typ.(IterableType).GetElementType(),
				),
			)
		case TypeString:
			return fmt.Sprintf("%s += %s", dest, g.codegenExpr(node.resolvedArgs["var"].expr, TypeString{}))
		}

	case "join":
		g.addImport("strings")
		listArg := node.resolvedArgs["list"]
		switch listArg.typ.(TypeSlice).GetElementType().(type) {
		case TypeString:
			callStr = fmt.Sprintf("strings.Join(%s, %s)",
				g.codegenExpr(node.resolvedArgs["list"].expr, NoCoercion{}),
				g.codegenExpr(node.resolvedArgs["sep"].expr, TypeString{}),
			)
		case TypeInt:
			g.addPreludeFunction("joinIntSlice")
			callStr = fmt.Sprintf("___joinIntSlice(%s, %s)",
				g.codegenVar(node.resolvedArgs["list"].expr.(*VarNode), NoCoercion{}),
				g.codegenExpr(node.resolvedArgs["sep"].expr, TypeString{}),
			)
		case TypeFloat:
			g.addPreludeFunction("joinFloatSlice")
			callStr = fmt.Sprintf("___joinFloatSlice(%s, %s)",
				g.codegenVar(node.resolvedArgs["list"].expr.(*VarNode), NoCoercion{}),
				g.codegenExpr(node.resolvedArgs["sep"].expr, TypeString{}),
			)
		default:
			panic("Joining not supported for this element type")
		}

	case "split":
		g.addImport("strings")
		callStr = fmt.Sprintf("strings.Split(%s, %s)",
			g.codegenExpr(node.resolvedArgs["string"].expr, TypeString{}),
			g.codegenExpr(node.resolvedArgs["sep"].expr, TypeString{}),
		)

	case "match":
		g.addPreludeFunction("regexMatch")
		callStr = fmt.Sprintf("___regexMatch(%s, %s)",
			g.codegenExpr(node.resolvedArgs["haystack"].expr, TypeString{}),
			g.codegenExpr(node.resolvedArgs["regex"].expr, TypeString{}),
		)

	case "capture":
		g.addPreludeFunction("regexCapture")
		callStr = fmt.Sprintf("___regexCapture(%s, %s)",
			g.codegenExpr(node.resolvedArgs["haystack"].expr, TypeString{}),
			g.codegenExpr(node.resolvedArgs["regex"].expr, TypeString{}),
		)

	case "find":
		g.addPreludeFunction("regexFind")
		callStr = fmt.Sprintf("___regexFind(%s, %s)",
			g.codegenExpr(node.resolvedArgs["haystack"].expr, TypeString{}),
			g.codegenExpr(node.resolvedArgs["regex"].expr, TypeString{}),
		)

	case "read":
		g.tmpVarCount++
		g.addImport("os")
		g.addImport("bufio")
		path := g.codegenExpr(node.resolvedArgs["path"].expr, TypeString{})
		genVar :=  g.codegenVar(&node.generatorVar, NoCoercion{})
		genVarSymbol, found := node.generatorBody.(*CompoundStatementNode).scope.lookupSymbol(genVar)
		if !found {
			panic("UNREACHABLE")
		}

		switch genVarSymbol.typ.(type) {
		case TypeString:
			g.addInitStatement(fmt.Sprintf("%s := ___scanner%d.Text()", genVar, g.tmpVarCount))
			g.addInitStatement(fmt.Sprintf("if !___chomp%d { %s=%s+\"\\n\" }", g.tmpVarCount, genVar, genVar))
		case TypeSlice:
			g.addImport("strings")
			g.addInitStatement(fmt.Sprintf("___string%d := ___scanner%d.Text()", g.tmpVarCount, g.tmpVarCount))
			g.addInitStatement(fmt.Sprintf("if !___chomp%d { ___string%d += \"\\n\" }", g.tmpVarCount, g.tmpVarCount))
			g.addInitStatement(fmt.Sprintf("%s := strings.Split(___string%d, %s)", genVar, g.tmpVarCount, g.codegenStringLiteral(node.resolvedArgs["sep"].expr.(*StringLiteralNode), NoCoercion{})))
		}

		g.addInitStatement(fmt.Sprintf("_ = %s", genVar))

		idxInitCode := ""
		if node.generatorHasIdx {
			genIdxVar := g.codegenVar(&node.generatorIdxVar, NoCoercion{})

			// This might generate slightly non-optimal go code, but is done to avoid
			// declaring the index variable outside the loop scope
			idxInitCode = fmt.Sprintf("___counter%d := -1", g.tmpVarCount)
			g.addInitStatement(fmt.Sprintf("___counter%d++", g.tmpVarCount))
			g.addInitStatement(fmt.Sprintf("%s := ___counter%d", genIdxVar, g.tmpVarCount))
		}

		body := g.codegenCompoundStatement(node.generatorBody.(*CompoundStatementNode))
		readCodeList := []string{
			fmt.Sprintf("___file%d, err := os.Open(%s)", g.tmpVarCount, path),
			g.indent("if err != nil {"),
			g.indent("    panic(\"Open fail\")"),
			g.indent("}"),
			g.indent(fmt.Sprintf("defer ___file%d.Close()", g.tmpVarCount)),
			g.indent(idxInitCode),
			g.indent(fmt.Sprintf("___scanner%d := bufio.NewScanner(___file%d)", g.tmpVarCount, g.tmpVarCount)),
			g.indent(fmt.Sprintf("___chomp%d := false", g.tmpVarCount)),
			g.indent(fmt.Sprintf("if %s { ___chomp%d = true }", g.codegenExpr(node.resolvedArgs["chomp"].expr, TypeBool{}), g.tmpVarCount)),
			g.indent(fmt.Sprintf("for ___scanner%d.Scan()", g.tmpVarCount)),
		}
		readCode := fmt.Sprintf("%s %s", strings.Join(readCodeList, "\n"), body)

		return readCode

	case "slurp":
		g.addPreludeFunction("slurpFile")
		callStr = fmt.Sprintf("___slurpFile(%s)",
			g.codegenExpr(node.resolvedArgs["path"].expr, TypeString{}),
		)

	default:
		panic("Unimplemented bulitin")
	}
	return g.coerce(callStr, builtin.returnType, coercion, CoercionModeDefault)
}

func (g *Generator) codegenReturn(node *ReturnNode) string {
	returnScope := g.scope.closestReturningScope()
	returnVal := g.codegenExpr(node.expr, returnScope.returnType)
	if returnScope.fallible {
		returnVal += ", nil"
	}
	return fmt.Sprintf("return %s", returnVal)
}

func (g *Generator) codegenFail(node *FailNode) string {
	// Find closest returnable scope
	returnScope := g.scope
	for returnScope.returnType.String() == "NoReturn" {
		returnScope = returnScope.parent
	}
	if !returnScope.fallible {
		panic("UNREACHABLE: Cannot fail from non-fallible function") // TODO: Actually check for this in typechecking!
	}

	failureString := g.codegenExpr(node.expr, TypeString{})

	if returnScope.returnType == (TypeVoid{}) {
		return fmt.Sprintf("return errors.New(%s)", failureString)
	}

	nilReturn := g.nilValue(returnScope.returnType)
	g.addImport("errors")
	return fmt.Sprintf("return %s,  errors.New(%s)", nilReturn, failureString)
}

func (g *Generator) codegenIf(node *IfNode) string {
	coerceType := node.compType
	if node.comp.Type() == VarNodeType || node.comp.Type() == AssignNodeType {
		coerceType = TypeBool{}
	}

	var elseCode string
	switch node.elseBody.Type() {
	case CompoundStatementNodeType:
		elseCode = " else " + g.codegenCompoundStatement(node.elseBody.(*CompoundStatementNode))
	case IfNodeType:
		elseCode = " else " + g.codegenIf(node.elseBody.(*IfNode))
	}

	comparison := g.codegenExpr(node.comp, coerceType)

	// Assignments in if statements generage prestatements that need to be added before the if statement
	prestatements := strings.Join(g.preStatements, "\n")
	g.preStatements = nil

	body := g.codegenCompoundStatement(node.body.(*CompoundStatementNode))
	return fmt.Sprintf(
		"%s\n%s %s %s%s",
		prestatements,
		g.indent("if"),
		comparison,
		body,
		elseCode,
	)
}

func (g *Generator) codegenForeach(node *ForeachNode) string {

	// Foreach loop with range: `for 1..10 -> x`
	if node.iterator.Type() == RangeNodeType {
		r := node.iterator.(*RangeNode)
		return fmt.Sprintf("for %s := %s; %s <= %s; %s++ %s",
			node.variable.token.str,
			g.codegenExpr(r.from, TypeInt{}),
			node.variable.token.str,
			g.codegenExpr(r.to, TypeInt{}),
			node.variable.token.str,
			g.codegenCompoundStatement(node.body.(*CompoundStatementNode)),
		)
	} else { // Foreach loop with iterator: `for list -> x`
		idxVarName := "_"
		if node.hasIdx {
			idxVarName = node.idxVariable.token.str
		}
		return fmt.Sprintf("for %s, %s := range %s %s",
			idxVarName,
			g.codegenVar(&node.variable, NoCoercion{}),
			g.codegenExpr(node.iterator, NoCoercion{}),
			g.codegenCompoundStatement(node.body.(*CompoundStatementNode)),
		)
	}
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
		return g.codegenFunctionCall(node.(*FunctionCallNode), NoCoercion{})
	case ReturnNodeType:
		return g.codegenReturn(node.(*ReturnNode))
	case FailNodeType:
		return g.codegenFail(node.(*FailNode))
	case IfNodeType:
		return g.codegenIf(node.(*IfNode))
	case ForeachNodeType:
		return g.codegenForeach(node.(*ForeachNode))
	case ContinueNodeType:
		return "continue"
	case BreakNodeType:
		return "break"
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
	case IndexedVarNodeType:
		return g.codegenIndexedVar(node.(*IndexedVarNode), coercion) //.(TypeSlice).ElementType)
	case FunctionCallNodeType:
		return g.codegenFunctionCall(node.(*FunctionCallNode), coercion)
	case SliceLiteralNodeType:
		return g.codegenSliceLiteral(node.(*SliceLiteralNode), coercion)
	case RangeNodeType:
		g.addPreludeFunction("createRange")
		return fmt.Sprintf("___createRange(%s, %s)", g.codegenExpr(node.(*RangeNode).from, TypeInt{}), g.codegenExpr(node.(*RangeNode).to, TypeInt{}))
	case AssignNodeType:
		return g.codegenAssignExpr(node.(*AssignNode), coercion)
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
	generator := Generator{0, nil, []string{}, make(map[string]bool), make(map[string]bool), []string{}, []string{}, []string{}, []string{}, 0, false, 0}
	code := generator.codegenProgram(root)

	if len(generator.errors) > 0 {
		return "", fmt.Errorf(strings.Join(generator.errors, "\n"))
	}

	return code, nil
}
