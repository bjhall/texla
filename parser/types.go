package parser

type Type interface {
	String() string
}

type IterableType interface {
	Type
	GetElementType() Type
}

type TypeInt struct{}

func (t TypeInt) String() string { return "int" }

type TypeFloat struct{}

func (t TypeFloat) String() string { return "float" }

type TypeString struct{}

func (t TypeString) String() string       { return "str" }
func (t TypeString) GetElementType() Type { return TypeString{} }

type TypeBool struct{}

func (t TypeBool) String() string { return "bool" }

type TypeUndetermined struct{}

func (t TypeUndetermined) String() string { return "Undetermined" }

type NoCoercion struct{}

func (t NoCoercion) String() string { return "NoCoercion" }

type NoReturn struct{}

func (t NoReturn) String() string { return "NoReturn" }

type TypeVoid struct{}

func (t TypeVoid) String() string { return "Void" }

type TypeSlice struct {
	ElementType Type
}

func (t TypeSlice) String() string       { return "[]" + t.ElementType.String() }
func (t TypeSlice) GetElementType() Type { return t.ElementType }

type TypeGenerator struct {
	ElementType Type
}
func (t TypeGenerator) String() string       { return "Generator<" + t.ElementType.String() + ">" }
func (t TypeGenerator) GetElementType() Type { return t.ElementType }


func isAppendable(t Type) bool {
	switch t.(type) {
	case TypeSlice, TypeString:
		return true
	default:
		return false
	}
}

func isGeneric(t Type) bool {
	switch t.(type) {
	case TypeInt, TypeFloat, TypeString, TypeBool, TypeUndetermined, TypeVoid, NoCoercion, NoReturn, TypeSlice, TypeGenerator:
		return false
	default:
		return true
	}
}

// Generic types
type TypeAny struct{}

func (t TypeAny) String() string { return "any" }

type TypeAppendable struct{}

func (t TypeAppendable) String() string { return "appendable" }

//type MapType struct {
//	KeyType   Type
//	ValueType Type
//}
//func (t MapType) String() string { return "map[" + t.KeyType.String() + "]" + t.ValueType.String() }
