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

	default:
		panic(fmt.Sprintf("Typechecking not implemented for builtin %q", builtin.name))
	}
	return returnType
}

func (tc *TypeChecker) typecheckExpr(node Node) Type {
	switch node.Type() {
	case NumNodeType:
		return node.(*NumNode).NumType()
	case StringLiteralNodeType:
		return TypeString{}
	case BoolNodeType:
		return TypeBool{}
	case SliceLiteralNodeType:
		typeCoercionPrecedence := map[Type]int{TypeFloat{}: 3, TypeString{}: 2, TypeInt{}: 1}

		var coercionType Type
		var highestPrecedence int
		for _, elem := range node.(*SliceLiteralNode).elements {
			typ := tc.typecheckExpr(elem)
			precedence, found := typeCoercionPrecedence[typ]
			if !found {
				tc.error(fmt.Sprintf("Type %s not allowed in slice literal", typ))
				continue
			}
			if precedence > highestPrecedence {
				highestPrecedence = precedence
				coercionType = typ
			}
		}

		node.(*SliceLiteralNode).elementType = coercionType
		return TypeSlice{ElementType: coercionType}

	case BinOpNodeType:
		leftType := tc.typecheckExpr(node.(*BinOpNode).left)
		rightType := tc.typecheckExpr(node.(*BinOpNode).right)
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

	case UnaryOpNodeType:
		return tc.typecheckExpr(node.(*UnaryOpNode).expr)

	case VarNodeType:
		varSymbol, found := tc.scope.lookupSymbol(node.(*VarNode).token.str)
		if found {
			return varSymbol.typ
		}
		fmt.Println("UNREACHABLE: Trying to look up type of undefined variable")
		os.Exit(1)

	case IndexedVarNodeType:
		varSymbol, found := tc.scope.lookupSymbol(node.(*IndexedVarNode).token.str)
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

	case FunctionCallNodeType:
		fnNode := node.(*FunctionCallNode)
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

	case RangeNodeType:
		fromType := tc.typecheckExpr(node.(*RangeNode).from)
		if (fromType != TypeInt{}) {
			tc.error("The from value of a range must be integer")
		}
		toType := tc.typecheckExpr(node.(*RangeNode).from)
		if (toType != TypeInt{}) {
			tc.error("The to value of a range must be integer")
		}
		return TypeSlice{ElementType: TypeInt{}}
	default:
		fmt.Println("TODO: Typechecking not implemented for", node.Type())
		os.Exit(1)
	}
	return TypeUndetermined{}
}

func (tc *TypeChecker) traverse(node Node) {

	switch node.Type() {

	case ProgramNodeType:
		for _, function := range node.(*ProgramNode).functions {
			tc.traverse(function)
		}

	case FunctionNodeType:
		tc.traverse(node.(*FunctionNode).body)

	case CompoundStatementNodeType:
		tc.scope = node.(*CompoundStatementNode).scope
		for _, child := range node.(*CompoundStatementNode).children {
			tc.traverse(child)
		}
		tc.scope = tc.scope.parent

	case IfNodeType:
		compType := tc.typecheckExpr(node.(*IfNode).comp)
		node.(*IfNode).setCompType(compType)
		// TODO: Ensure that comparison is a boolean value
		tc.traverse(node.(*IfNode).comp)
		tc.traverse(node.(*IfNode).body)
		tc.traverse(node.(*IfNode).elseBody)

	case AssignNodeType:
		lhs := node.(*AssignNode).left
		rhs := node.(*AssignNode).right

		lhsSymbol, found := tc.scope.lookupSymbol(lhs.(*VarNode).token.str)

		if found && lhsSymbol.typ.String() == "Undetermined" {
			rhsType := tc.typecheckExpr(rhs)
			tc.scope.setSymbolType(lhs.(*VarNode).token.str, rhsType)
		} else {
			// Do nothing?
		}

		tc.traverse(rhs)

	case ReturnNodeType:
		node.(*ReturnNode).setType(tc.typecheckExpr(node.(*ReturnNode).expr))

	case FunctionCallNodeType:
		fnNode := node.(*FunctionCallNode)
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

		for _, argNode := range fnNode.resolvedArgs {
			tc.traverse(&argNode)
		}

	case IndexedVarNodeType:
		tc.traverse(node.(*IndexedVarNode).index)

	case ArgumentNodeType:
		tc.traverse(node.(*ArgumentNode).expr)

	case BinOpNodeType:
		tc.traverse(node.(*BinOpNode).left)
		tc.traverse(node.(*BinOpNode).right)

	case SliceLiteralNodeType:
		for _, el := range node.(*SliceLiteralNode).elements {
			tc.traverse(el)
		}

	case ForeachNodeType:
		var controlVarType Type

		switch node.(*ForeachNode).iterator.Type() {
		case RangeNodeType:
			controlVarType = TypeInt{}
		default:
			iterType := tc.typecheckExpr(node.(*ForeachNode).iterator)
			controlVarType = iterType.(IterableType).GetElementType()
		}
		node.(*ForeachNode).body.(*CompoundStatementNode).SetVarType(node.(*ForeachNode).variable.token.str, controlVarType)
		tc.traverse(node.(*ForeachNode).body)

	case StringLiteralNodeType, NumNodeType, BoolNodeType, VarNodeType, NoOpNodeType, UnaryOpNodeType, RangeNodeType:
		return

	default:
		fmt.Println("TYPECHECKING TODO:", node.Type())
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
