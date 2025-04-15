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
func (t TypeString) String() string { return "str" }
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


type TypeSlice struct{
	ElementType Type
}
func (t TypeSlice) String() string { return "[]" + t.ElementType.String() }
func (t TypeSlice) GetElementType() Type { return t.ElementType }

//type MapType struct {
//	KeyType   Type
//	ValueType Type
//}
//func (t MapType) String() string { return "map[" + t.KeyType.String() + "]" + t.ValueType.String() }
