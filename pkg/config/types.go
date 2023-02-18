package config

import "reflect"

// field maintains information about the struct field
type field struct {
	Name    string
	RfValue reflect.Value
	RfTags  reflect.StructTag
	Value   string
}
