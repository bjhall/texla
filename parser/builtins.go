package parser

type BuiltinFunc struct {
	name       string
	returnType Type
	parameters []ParameterNode
	generator  bool
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
	"add": {
		name:       "add",
		returnType: TypeVoid{},
		parameters: []ParameterNode{
			ParameterNode{name: "dest", typ: TypeSet{}},
			ParameterNode{name: "var", typ: TypeAny{}},
		},
	},
	"has": {
		name:       "has",
		returnType: TypeBool{},
		parameters: []ParameterNode{
			ParameterNode{name: "haystack", typ: TypeSet{}},
			ParameterNode{name: "needle", typ: TypeAny{}},
		},
	},
	"del": {
		name:       "del",
		returnType: TypeBool{},
		parameters: []ParameterNode{
			ParameterNode{name: "set", typ: TypeSet{}},
			ParameterNode{name: "value", typ: TypeAny{}},
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
	"split": {
		name:       "split",
		returnType: TypeSlice{ElementType: TypeString{}},
		parameters: []ParameterNode{
			ParameterNode{name: "string", typ: TypeString{}},
			ParameterNode{name: "sep", typ: TypeString{}},
		},
	},
	"read": {
		name: "read",
		returnType: TypeGenerator{ElementType: TypeUndetermined{}},
		generator: true,
		parameters: []ParameterNode{
			ParameterNode{name: "path", typ: TypeString{}},
			ParameterNode{name: "chomp", typ: TypeBool{}, hasDefault: true, defaultValue: "true"},
			ParameterNode{name: "sep", typ: TypeString{}, hasDefault: true, defaultValue: ""},
		},
	},
	"slurp": {
		name: "slurp",
		returnType: TypeString{},
		parameters: []ParameterNode{
			ParameterNode{name: "path", typ: TypeString{}},
		},
	},
	"match": {
		name: "match",
		returnType: TypeBool{},
		parameters: []ParameterNode{
			ParameterNode{name: "haystack", typ: TypeString{}},
			ParameterNode{name: "regex", typ: TypeString{}},
		},
	},
	"capture": {
		name: "capture",
		returnType: TypeSlice{ElementType: TypeString{}},
		parameters: []ParameterNode{
			ParameterNode{name: "haystack", typ: TypeString{}},
			ParameterNode{name: "regex", typ: TypeString{}},
		},
	},
	"find": {
		name: "find",
		returnType: TypeSlice{ElementType: TypeString{}},
		parameters: []ParameterNode{
			ParameterNode{name: "haystack", typ: TypeString{}},
			ParameterNode{name: "regex", typ: TypeString{}},
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
