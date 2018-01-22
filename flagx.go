package flagx

import (
	"encoding"
	"reflect"
)

// Value is like flag.Getter (which is a superset of flag.Value)
type Value interface {
	String() string
	Set(string) error
	Get() interface{}
}

// Dummy is a flag.Value that does nothing.
type Dummy struct{}

func (Dummy) String() string { return "" }

func (Dummy) Set(s string) error { return nil }

func (Dummy) Get() interface{} { return nil }

// stringSetter is the subset of flag.Value for setting a value from a string
type stringSetter interface {
	// See flag.Value
	Set(string) error
}

var (
	textUnmarshalerType = reflect.TypeOf(new(encoding.TextUnmarshaler)).Elem()
	stringSetterType    = reflect.TypeOf(new(stringSetter)).Elem()
)

func setterFor(typ reflect.Type) func(target reflect.Value, value string) error {
	switch {
	case reflect.PtrTo(typ).Implements(stringSetterType):
		return func(target reflect.Value, value string) error {
			return target.Addr().Interface().(stringSetter).Set(value)
		}
	case reflect.PtrTo(typ).Implements(textUnmarshalerType):
		return func(target reflect.Value, value string) error {
			return target.Addr().Interface().(encoding.TextUnmarshaler).UnmarshalText([]byte(value))
		}
	case typ.Implements(stringSetterType):
		return func(target reflect.Value, value string) error {
			return target.Interface().(stringSetter).Set(value)
		}
	case typ.Implements(textUnmarshalerType):
		return func(target reflect.Value, value string) error {
			return target.Interface().(encoding.TextUnmarshaler).UnmarshalText([]byte(value))
		}
	case typ.Kind() == reflect.String:
		return func(target reflect.Value, value string) error {
			target.SetString(value)
			return nil
		}
	default:
		return nil
	}
}
