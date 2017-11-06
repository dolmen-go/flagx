package flagx

import (
	"fmt"
	"reflect"
	"strings"
)

func Map(m interface{}, parseValue func(string) (interface{}, error)) Value {
	v := reflect.ValueOf(m)
	if v.IsNil() || v.Kind() != reflect.Map {
		panic("non-nil pointer to a map expected")
	}
	// check that keys are strings
	if v.Type().Key().Kind() != reflect.String {
		panic("keys must be of type string")
	}
	return &stringMap{v, parseValue}
}

type stringMap struct {
	Map   reflect.Value
	Parse func(string) (interface{}, error)
}

func (m *stringMap) String() string {
	return ""
}

func (m *stringMap) Set(s string) error {
	i := strings.IndexByte(s, '=')
	if i < 0 {
		return fmt.Errorf("%q: '=' expected")
	}
	key := s[:i]
	value, err := m.Parse(s[i+1:])
	if err != nil {
		return err
	}
	m.Map.SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(value))
	return nil
}

func (m *stringMap) Get() interface{} {
	return m.Map.Interface()
}
