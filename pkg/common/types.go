package common

import "reflect"

// Field maintains information about the struct field
type Field struct {
	Name    string
	RfValue reflect.Value
	RfTags  reflect.StructTag
	Value   string
}
