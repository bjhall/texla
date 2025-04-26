package parser

type BuiltinFunc struct {
	name       string
	returnType Type
	parameters []ParameterNode
}

var builtins = map[string]BuiltinFunc{
	"len": {
		name:       "len",
		returnType: TypeInt{},
		parameters: []ParameterNode{
			ParameterNode{name: "var", typ: TypeString{}},
		},
	},
	"append": {
		name:       "append",
		returnType: TypeVoid{},
		parameters: []ParameterNode{
			ParameterNode{name: "dest", typ: TypeAppendable{}},
			ParameterNode{name: "var", typ: TypeAny{}},
		},
	},
	"join": {
		name:       "join",
		returnType: TypeString{},
		parameters: []ParameterNode{
			ParameterNode{name: "list", typ: TypeSlice{}},
			ParameterNode{name: "sep", typ: TypeString{}},
		},
	},
}

func isBuiltin(name string) bool {
	_, found := builtins[name]
	return found
}

func (b *BuiltinFunc) getParamType(paramName string) Type {
	for _, param := range b.parameters {
		if param.name == paramName {
			return param.typ
		}
	}
	panic("Invalid parameter to builtin")
}
