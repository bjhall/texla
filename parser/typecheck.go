package parser

import (
	"fmt"
	"os"
	"strings"
)

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
		return TypeString
	case BinOpNodeType:
		leftType := tc.typecheckExpr(node.(*BinOpNode).left)
		rightType := tc.typecheckExpr(node.(*BinOpNode).right)
		if leftType == rightType {
			return leftType
		}
		if leftType == TypeFloat || rightType == TypeFloat {
			return TypeFloat
		}
		if leftType == TypeString || rightType == TypeString {
			return TypeFloat
		}
		return TypeInt
	case VarNodeType:
		varSymbol, found :=  tc.scope.lookupSymbol(node.(*VarNode).token.str)
		if found {
			return varSymbol.typ
		}
		fmt.Println("UNREACHABLE: Trying to look up type of undefined variable")
		os.Exit(1)
	case FunctionCallNodeType:
		funcSymbol, found := tc.scope.lookupSymbol(node.(*FunctionCallNode).name)
		if found {
			return funcSymbol.typ
		}
		fmt.Println("UNREACHABLE: Trying to loop up type of non-existing function")
		os.Exit(1)
	default:
		node.Print(1)
		fmt.Println("TODO: Typechecking not implemented for", node.Type())
		os.Exit(1)
	}
	return TypeUndetermined
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

	case AssignNodeType:
		lhs := node.(*AssignNode).left
		rhs := node.(*AssignNode).right

		lhsSymbol, found := tc.scope.lookupSymbol(lhs.(*VarNode).token.str)

		if found && lhsSymbol.typ == TypeUndetermined {
			rhsType := tc.typecheckExpr(rhs)
			tc.scope.setSymbolType(lhs.(*VarNode).token.str, rhsType)
		} else {
			panic("TODO: Variable reassignment not implemented")
		}

	case ReturnNodeType:
		node.(*ReturnNode).setType(tc.typecheckExpr(node.(*ReturnNode).expr))

	case FunctionCallNodeType:
		functionName := node.(*FunctionCallNode).name

		symbol, found := tc.scope.lookupSymbol(functionName)

		// Check if symbol is declared
		if !found {
			if functionName == "print" {
				// TODO: HANDLE BUILTINS
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

		// Check if the right number of arguments was provided
		numArguments := len(node.(*FunctionCallNode).arguments)
		if numArguments != len(symbol.parameterTypes) {
			tc.error(fmt.Sprintf("Wrong number of arguments to %s, expected %d, got %d", functionName, len(symbol.parameterTypes), numArguments))
		}


		// Infer types of function call argumemnts
		for i, argument := range node.(*FunctionCallNode).arguments {
			expectedType := symbol.parameterTypes[i]
			givenType := tc.typecheckExpr(argument)
			if givenType == TypeFloat && expectedType == TypeInt && argument.Type() == NumNodeType {
				tc.addImport("math")
			}
			node.(*FunctionCallNode).argumentTypes = append(node.(*FunctionCallNode).argumentTypes, givenType)
		}

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
