package flagx

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type IntSlice struct {
	Slice *[]int
}

func (is IntSlice) String() string {
	if is.Slice == nil {
		// When called by flag.isZeroValue
		return ""
	}
	if *is.Slice == nil {
		return ""
	}
	str := make([]string, len(*is.Slice))
	for i := range str {
		str[i] = strconv.Itoa((*is.Slice)[i])
	}
	return strings.Join(str, ",")
}

func (is IntSlice) Set(s string) (err error) {
	str := strings.Split(s, ",")
	for _, s := range str {
		s = strings.TrimSpace(s)
		n, err := strconv.ParseInt(s, 0, 0)
		if err != nil {
			return fmt.Errorf("%q: %s", s, err)
		}
		*is.Slice = append(*is.Slice, int(n))
	}
	return nil
}

func (is IntSlice) Get() interface{} {
	return *is.Slice
}

// Slice wraps any pointer to any slice type to expose it
// as a flag.Value
// Uses reflect.
func Slice(sl interface{}, separator string, parse func(string) (interface{}, error)) Value {
	v := reflect.ValueOf(sl)
	if v.Type().Kind() != reflect.Ptr {
		panic("pointer expected")
	}
	if v.IsNil() {
		panic("non-nil pointer expected")
	}
	if v.Elem().Kind() != reflect.Slice {
		panic("non-nil pointer to a slice expected")
	}
	return &slice{v, separator, parse}
}

type slice struct {
	Slice     reflect.Value
	Separator string
	Parse     func(string) (interface{}, error)
}

func (sl *slice) String() string {
	return ""
}

func (sl *slice) appnd(s string) error {
	v, err := sl.Parse(s)
	if err != nil {
		return err
	}
	sl.Slice.Elem().Set(reflect.Append(sl.Slice.Elem(), reflect.ValueOf(v)))
	return nil
}

func (sl *slice) Set(s string) error {
	if sl.Separator != "" {
		for _, item := range strings.Split(s, sl.Separator) {
			if err := sl.appnd(item); err != nil {
				return err
			}
		}
		return nil
	}
	return sl.appnd(s)
}

func (sl *slice) Get() interface{} {
	return sl.Slice.Elem().Interface()
}
