package flagx

import (
	"fmt"
	"reflect"
	"strings"
)

// Map wraps any map to expose it as a flag.Value.
//
// The value type of the map may implement flag.Value. In that case the
// Set() method will be called on the target element.
// The element type of the map may implement encoding.TextUnmarshaler. In
// that case the UnmarshalText() method will be called on the target element.
//
// The parse func is optional. If it is set, it must return a value
// assignable to an element of the slice. If the returned value is a bare
// string, it will pass through Set() or UnmarshalText() if the type implements
// it (see above).
func Map(m interface{}, parseValue func(string) (interface{}, error)) Value {
	v := reflect.ValueOf(m)
	if v.IsNil() || v.Kind() != reflect.Map {
		panic("non-nil pointer to a map expected")
	}
	// check that keys are strings
	// FIXME check must be stricter
	if v.Type().Key().Kind() != reflect.String {
		panic("keys must be of type string")
	}

	valueType := v.Type().Elem()
	set := setterFor(valueType)
	var buildValue func(value string) (reflect.Value, error)
	if parseValue == nil {
		if set == nil {
			panic(fmt.Errorf("invalid slice type: %s doesn't implement encoding.TextUnmarshaler or flag.Value", v.Type()))
		}
		if valueType.Kind() == reflect.Interface {
			panic("a parse function must be provided to build a concrete value")
		}
		buildValue = func(value string) (reflect.Value, error) {
			v := reflect.New(valueType).Elem()
			return v, set(v, value)
		}
	} else {
		if set == nil {
			buildValue = func(value string) (reflect.Value, error) {
				v, err := parseValue(value)
				if err != nil {
					return reflect.Value{}, err
				}
				return reflect.ValueOf(&v).Elem().Elem(), nil
			}
		} else {
			buildValue = func(value string) (reflect.Value, error) {
				v, err := parseValue(value)
				if err != nil {
					return reflect.Value{}, err
				}
				if s, isString := v.(string); isString {
					rv := reflect.New(valueType).Elem()
					return rv, set(rv, s)
				}
				return reflect.ValueOf(&v).Elem().Elem(), nil
			}
		}
	}

	return &stringMap{v, buildValue}
}

type stringMap struct {
	Map        reflect.Value
	buildValue func(string) (reflect.Value, error)
}

func (m *stringMap) String() string {
	return ""
}

func (m *stringMap) Set(s string) error {
	i := strings.IndexByte(s, '=')
	if i < 0 {
		return fmt.Errorf("%q: '=' expected", s)
	}
	key := s[:i]

	value, err := m.buildValue(s[i+1:])
	if err != nil {
		return err
	}

	m.Map.SetMapIndex(reflect.ValueOf(key), value)
	return nil
}

func (m *stringMap) Get() interface{} {
	return m.Map.Interface()
}
