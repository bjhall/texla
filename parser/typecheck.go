package parser

import (
	"fmt"
	"os"
	"strings"
)

const ArgNotProvided = -1

type TypeChecker struct {
	scope       *Scope
	errors      []string
	imports     map[string]bool
}

func (tc *TypeChecker) error(errorStr string) {
	tc.errors = append(tc.errors, fmt.Sprintf("%s", errorStr))
}


func (tc *TypeChecker) addImport(name string) {
	tc.imports[name] = true
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
			fmt.Println("%s is not indexable", t)
		}

	case FunctionCallNodeType:
		funcSymbol, found := tc.scope.lookupSymbol(node.(*FunctionCallNode).name)
		if found {
			return funcSymbol.typ
		}
		fmt.Println("UNREACHABLE: Trying to look up type of non-existing function")
		os.Exit(1)
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
			panic("TODO: Variable reassignment not implemented")
		}

		tc.traverse(rhs)

	case ReturnNodeType:
		node.(*ReturnNode).setType(tc.typecheckExpr(node.(*ReturnNode).expr))

	case FunctionCallNodeType:
		functionName := node.(*FunctionCallNode).name

		symbol, found := tc.scope.lookupSymbol(functionName)

		// Check if symbol is declared
		if !found {
			if functionName == "print" {
				// TODO: HANDLE BUILTINS
				for _, arg := range node.(*FunctionCallNode).arguments {
					tc.traverse(arg.(*ArgumentNode).expr)
				}
				return
			} else {
				tc.error(fmt.Sprintf("No function named %q exists in current scope", functionName))
				return
			}
		}

		// Check if the symbol is a function
		if symbol.category != FunctionSymbol {
			tc.error(fmt.Sprintf("%q is not a function", functionName))
		}

		// Validate and typecheck arguments
		for i, param := range symbol.paramsNode.parameters {
			var givenType Type
			var argIdx int
			found := false
			for j, n := range node.(*FunctionCallNode).arguments {
				arg := n.(*ArgumentNode)
				if (arg.named && arg.paramName == param.name) || (!arg.named && arg.order == i) {
					argIdx = j
					givenType = tc.typecheckExpr(arg.expr)
					found = true
					break
				}
			}
			if !found {
				if !param.hasDefault {
					tc.error(fmt.Sprintf("Value missing for argument %q (%s) of function %q", param.name, param.typ, functionName))
					return
				}
				node.(*FunctionCallNode).appendArgumentOrder(ArgNotProvided)
			} else {
				node.(*FunctionCallNode).arguments[argIdx].(*ArgumentNode).typ = givenType
				node.(*FunctionCallNode).appendArgumentOrder(argIdx)
				tc.traverse(node.(*FunctionCallNode).arguments[argIdx])
			}
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

	case StringLiteralNodeType, NumNodeType, BoolNodeType, VarNodeType, NoOpNodeType, UnaryOpNodeType:
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
