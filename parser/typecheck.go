package parser

import (
	"fmt"
	"os"
	"strings"
)

type TypeChecker struct {
	scope   *Scope
	errors  []string
	imports map[string]bool
}

func (tc *TypeChecker) error(errorStr string) {
	tc.errors = append(tc.errors, fmt.Sprintf("%s", errorStr))
}

func (tc *TypeChecker) addImport(name string) {
	tc.imports[name] = true
}

func (tc *TypeChecker) typecheckBuiltin(node Node) Type {
	var returnType Type
	fnNode := node.(*FunctionCallNode)
	builtin := builtins[fnNode.name]
	if !isGeneric(builtin.returnType) {
		returnType = builtin.returnType
	} else {
		panic("Resolve generic return type")
	}

	err := fnNode.matchArgsToParams(builtin.parameters)
	if err != nil {
		tc.error(err.Error())
	}

	switch builtin.name {
	case "append":
		containerType := tc.typecheckExpr(fnNode.resolvedArgs["dest"].expr)
		if !isAppendable(containerType) {
			tc.error(fmt.Sprintf("append() cannot be used on type %q", containerType))
		}
		node.(*FunctionCallNode).setArgType("dest", containerType)

	case "add":
		containerType := tc.typecheckExpr(fnNode.resolvedArgs["dest"].expr)
		if !isSettable(containerType) {
			tc.error(fmt.Sprintf("add() can only be used on sets, not %q", containerType))
		}
		node.(*FunctionCallNode).setArgType("dest", containerType)

	case "has":
		containerType := tc.typecheckExpr(fnNode.resolvedArgs["haystack"].expr)
		// TODO: Allow contains to be used on slices and strings too?
		if !isSettable(containerType) {
			tc.error(fmt.Sprintf("contains() can only be used on sets, not %q", containerType))
		}
		node.(*FunctionCallNode).setArgType("haystack", containerType)

	case "del":
		containerType := tc.typecheckExpr(fnNode.resolvedArgs["set"].expr)
		if !isSettable(containerType) {
			tc.error(fmt.Sprintf("del() can only be used on sets, not %q", containerType))
		}
		node.(*FunctionCallNode).setArgType("set", containerType)

	case "union":
		set1Type := tc.typecheckExpr(fnNode.resolvedArgs["set1"].expr)
		set2Type := tc.typecheckExpr(fnNode.resolvedArgs["set2"].expr)
		if !isSettable(set1Type) {
			tc.error(fmt.Sprintf("union() can only be used on sets, not %q", set1Type))
		}
		if !isSettable(set2Type) {
			tc.error(fmt.Sprintf("union() can only be used on sets, not %q", set2Type))
		}
		if set1Type.(TypeSet).ElementType != set2Type.(TypeSet).ElementType {
			tc.error(fmt.Sprintf("union() can only be used on sets with the same element type"))
		}
		//fmt.Println(set1Type.(TypeSet).ElementType, set2Type.(TypeSet).ElementType)
		node.(*FunctionCallNode).setArgType("set1", set1Type)
		node.(*FunctionCallNode).setArgType("set2", set1Type)
		returnType = set1Type

	case "to_set":
		containerType := tc.typecheckExpr(fnNode.resolvedArgs["slice"].expr)
		if _, isSlice := containerType.(TypeSlice); !isSlice {
			tc.error(fmt.Sprintf("to_setn() can only be used on slices", containerType))
		}
		returnType = TypeSet{ElementType: containerType.(TypeSlice).ElementType}

	case "len":
		containerType := tc.typecheckExpr(fnNode.resolvedArgs["var"].expr)
		if !isAppendable(containerType) {
			tc.error(fmt.Sprintf("len() cannot be used on type %q", containerType))
		}
	case "join":
		containerType := tc.typecheckExpr(fnNode.resolvedArgs["list"].expr)
		if _, isSlice := containerType.(TypeSlice); !isSlice {
			tc.error(fmt.Sprintf("join() can only be used on slices", containerType))
		}
		node.(*FunctionCallNode).setArgType("list", containerType)
	case "read":
		separator, sepIsString :=fnNode.resolvedArgs["sep"].expr.(*StringLiteralNode)
		if !sepIsString {
			tc.error(fmt.Sprintf("sep argument for read() must be string literal"))
		}
		if separator.token.str == "" {
			returnType = TypeString{}
		} else {
			returnType = TypeSlice{ElementType: TypeString{}}
		}
	case "split", "match", "capture", "find", "slurp":
		// Do nothing?
	default:
		panic(fmt.Sprintf("Typechecking not implemented for builtin %q", builtin.name))
	}
	return returnType
}

func (tc *TypeChecker) typecheckExprList(nodes []Node) Type {
	typeCoercionPrecedence := map[Type]int{TypeString{}: 3, TypeFloat{}: 2, TypeInt{}: 1}
	var coercionType Type
	var highestPrecedence int
	for _, elem := range nodes {
		typ := tc.typecheckExpr(elem)
		precedence, found := typeCoercionPrecedence[typ]
		if !found {
			tc.error(fmt.Sprintf("Type %s not allowed in expression list", typ))
			continue
		}
		if precedence > highestPrecedence {
			highestPrecedence = precedence
			coercionType = typ
		}
	}
	return coercionType
}

func (tc *TypeChecker) typecheckExpr(node Node) Type {
	switch n := node.(type) {
	case *NumNode:
		return n.NumType()
	case *StringLiteralNode:
		return TypeString{}
	case *BoolNode:
		return TypeBool{}
	case *SliceLiteralNode:
		elementType := tc.typecheckExprList(n.elements)
		n.elementType = elementType
		return TypeSlice{ElementType: elementType}
	case *SetLiteralNode:
		elementType := tc.typecheckExprList(n.elements)
		n.elementType = elementType
		return TypeSet{ElementType: elementType}
	case *BinOpNode:
		leftType := tc.typecheckExpr(n.left)
		rightType := tc.typecheckExpr(n.right)
		if leftType == rightType {
			return leftType
		}
		if leftType.String() == "float" || rightType.String() == "float" {
			return TypeFloat{}
		}
		if leftType.String() == "string" || rightType.String() == "string" {
			return TypeFloat{}
		}
		return TypeInt{}

	case *UnaryOpNode:
		return tc.typecheckExpr(n.expr)

	case *VarNode:
		varSymbol, found := tc.scope.lookupSymbol(n.token.str)
		if found {
			return varSymbol.typ
		}
		fmt.Println("UNREACHABLE: Trying to look up type of undefined variable")
		os.Exit(1)

	case *IndexedVarNode:
		varSymbol, found := tc.scope.lookupSymbol(n.token.str)
		if !found {
			fmt.Println("UNREACHABLE: Trying to look up type of undefined indexed variable")
			os.Exit(1)
		}
		switch t := varSymbol.typ.(type) {
		case TypeSlice:
			return t.ElementType
		case TypeString:
			return TypeString{}
		default:
			fmt.Printf("%s is not indexable\n", t)
		}

	case *FunctionCallNode:
		fnNode := n
		functionName := fnNode.name
		if isBuiltin(functionName) {
			return tc.typecheckBuiltin(fnNode)
		}

		funcSymbol, found := tc.scope.lookupSymbol(functionName)

		if found {
			err := fnNode.matchArgsToParams(funcSymbol.paramsNode.parameters)
			if err != nil {
				tc.error(err.Error())
			}
			return funcSymbol.typ
		}
		fmt.Println("UNREACHABLE: Trying to look up type of non-existing function")
		os.Exit(1)

	case *RangeNode:
		fromType := tc.typecheckExpr(n.from)
		if (fromType != TypeInt{}) {
			tc.error("The from value of a range must be integer")
		}
		toType := tc.typecheckExpr(n.from)
		if (toType != TypeInt{}) {
			tc.error("The to value of a range must be integer")
		}
		return TypeSlice{ElementType: TypeInt{}}
	case *AssignNode:
		if n.expression == false {
			panic("Typechecking non-exression assignment")
		}
		lhs := n.left
		lhsSymbol, found := tc.scope.lookupSymbol(lhs.(*VarNode).token.str)
		if found && lhsSymbol.typ.String() == "Undetermined" {
			rhsType := tc.typecheckExpr(n.right)
			tc.scope.setSymbolType(lhs.(*VarNode).token.str, rhsType)
		}
		return tc.typecheckExpr(n.right)

	default:
		fmt.Printf("TODO: Typechecking not implemented for: %T\n")
		os.Exit(1)
	}
	return TypeUndetermined{}
}

func (tc *TypeChecker) traverse(node Node) {

	switch n := node.(type) {

	case *ProgramNode:
		for _, function := range n.functions {
			tc.traverse(function)
		}

	case *FunctionNode:
		tc.traverse(n.body)

	case *CompoundStatementNode:
		tc.scope = n.scope
		for _, child := range n.children {
			tc.traverse(child)
		}
		tc.scope = tc.scope.parent

	case *IfNode:
		compType := tc.typecheckExpr(n.comp)
		node.(*IfNode).setCompType(compType)
		// TODO: Ensure that comparison is a boolean value
		tc.traverse(n.comp)
		tc.traverse(n.body)
		tc.traverse(n.elseBody)

	case *AssignNode:
		if !n.expression {
			lhsSymbol, found := tc.scope.lookupSymbol(n.left.(*VarNode).token.str)
			if found && lhsSymbol.typ.String() == "Undetermined" {
				rhsType := tc.typecheckExpr(n.right)
				tc.scope.setSymbolType(n.left.(*VarNode).token.str, rhsType)
			} else {
				// Do nothing?
			}
			tc.traverse(n.right)
		}

	case *ReturnNode:
		n.setType(tc.typecheckExpr(n.expr))

	case *FailNode:
		if !tc.scope.closestReturningScope().fallible {
			tc.error("Cannot use `fail` in non-fallible function")
		}
		exprType := tc.typecheckExpr(n.expr)
		if exprType != (TypeString{}) {
			tc.error("Failure expression must be a string")
		}


	case *FunctionCallNode:
		fnNode := n
		functionName := fnNode.name

		var parameters []ParameterNode

		if isBuiltin(functionName) {
			_ = tc.typecheckBuiltin(node)
			parameters = builtins[functionName].parameters
		} else {
			symbol, found := tc.scope.lookupSymbol(functionName)
			if found {
				if symbol.category != FunctionSymbol {
					tc.error(fmt.Sprintf("%q is not a function", functionName))
				}
				parameters = symbol.paramsNode.parameters

				// Check that errors are handled correctly
				if symbol.fallible && !fnNode.errorHandled {
					tc.error(fmt.Sprintf("Function %q can return an error, but it is not handled", functionName))
				}
				if !symbol.fallible && fnNode.errorHandled {
					tc.error(fmt.Sprintf("Function %q is not fallible, do not put ? after the call to it", functionName))
				}

			} else if functionName == "print" {
				// TODO: Make print a builtin
				for _, arg := range fnNode.arguments {
					tc.traverse(arg.(*ArgumentNode).expr)
				}
				return
			} else {
				tc.error(fmt.Sprintf("No function named %q exists in current scope", functionName))
				return
			}
		}

		err := fnNode.matchArgsToParams(parameters)
		if err != nil {
			tc.error(err.Error())
		}

		if fnNode.errorBody != nil {
			tc.traverse(fnNode.errorBody)
		}

		for _, argNode := range fnNode.resolvedArgs {
			tc.traverse(&argNode)
		}

		if fnNode.generatorVar != (VarNode{}) {
			//controlVarType := builtins[functionName].returnType.(TypeGenerator).GetElementType()
			controlVarType := tc.typecheckBuiltin(node)
			fnNode.generatorBody.(*CompoundStatementNode).SetVarType(fnNode.generatorVar.token.str, controlVarType)
		}
		if fnNode.generatorBody != nil {
			tc.traverse(fnNode.generatorBody)
		}


	case *IndexedVarNode:
		tc.traverse(n.index)

	case *RangeNode:
		tc.traverse(n.from)
		tc.traverse(n.to)

	case *ArgumentNode:
		tc.traverse(n.expr)

	case *BinOpNode:
		tc.traverse(n.left)
		tc.traverse(n.right)

	case *SliceLiteralNode:
		for _, el := range n.elements {
			tc.traverse(el)
		}

	case *SetLiteralNode:
		for _, el := range n.elements {
			tc.traverse(el)
		}

	case *ForeachNode:
		var controlVarType Type

		switch n.iterator.(type) {
		case *RangeNode:
			controlVarType = TypeInt{}
		default:
			iterType := tc.typecheckExpr(n.iterator)
			controlVarType = iterType.(IterableType).GetElementType()
		}
		n.body.(*CompoundStatementNode).SetVarType(n.variable.token.str, controlVarType)
		tc.traverse(n.body)

	case *IncNode:
		varSymbol, _ := tc.scope.lookupSymbol(n.varName)
		switch varSymbol.typ.(type) {
		case TypeInt, TypeFloat:
		default:
			tc.error(fmt.Sprintf("Cannot use ++ operator on non-numeric types"))
		}
	case *DecNode:
		varSymbol, _ := tc.scope.lookupSymbol(n.varName)
		switch varSymbol.typ.(type) {
		case TypeInt, TypeFloat:
		default:
			tc.error(fmt.Sprintf("Cannot use -- operator on non-numeric types"))
		}

	case *StringLiteralNode, *NumNode, *BoolNode, *VarNode, *NoOpNode, *UnaryOpNode, *ContinueNode, *BreakNode:
		return

	default:
		fmt.Printf("TYPECHECKING TODO: %T\n", node)
		os.Exit(1)

	}
}

func CheckTypes(root Node) (Node, error) {
	typeChecker := TypeChecker{nil, []string{}, make(map[string]bool)}

	typeChecker.traverse(root)

	for importName, _ := range typeChecker.imports {
		root.(*ProgramNode).addImport(importName)
	}

	if len(typeChecker.errors) > 0 {
		return root, fmt.Errorf(strings.Join(typeChecker.errors, "\n"))
	}

	return root, nil
}
