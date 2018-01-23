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

// Slice wraps any pointer to a slice to expose it as a flag.Value.
//
// The element type of the slice may implement flag.Value. In that case the
// Set() method will be called on the target element.
// The element type of the slice may implement encoding.TextUnmarshaler. In
// that case the UnmarshalText() method will be called on the target element.
//
// The parse func is optional. If it is set, it must return a value
// assignable to an element of the slice. If the returned value is a bare
// string, it will pass through Set() or UnmarshalText() if the type implements
// it (see above).
func Slice(sl interface{}, separator string, parse func(string) (interface{}, error)) Value {
	v := reflect.ValueOf(sl)
	if v.Type().Kind() != reflect.Ptr {
		panic("pointer expected")
	}
	if v.IsNil() {
		panic("non-nil pointer expected")
	}
	v = v.Elem() // v now wraps the slice
	if v.Kind() != reflect.Slice {
		panic("non-nil pointer to a slice expected")
	}
	itemType := v.Type().Elem()
	setter := setterFor(itemType)
	if parse == nil {
		if setter == nil {
			panic(fmt.Errorf("invalid slice type: %s doesn't implement encoding.TextUnmarshaler or flag.Value", v.Type()))
		}
		if itemType.Kind() == reflect.Interface {
			panic("a parse function must be provided to build a concrete value")
		}
	} else {
		if setter == nil {
			setter = func(target reflect.Value, value string) error {
				v, err := parse(value)
				if err != nil {
					return err
				}
				// FIXME This need more testing!
				if itemType.Kind() == reflect.Interface {
					if itemType.NumMethod() == 0 {
						target.Set(reflect.ValueOf(&v).Elem())
					} else {
						target.Set(reflect.ValueOf(&v).Elem().Elem().Convert(itemType))
					}
					return nil
				}
				target.Set(reflect.ValueOf(&v).Elem().Elem())
				return nil
			}
		} else {
			setString := setter
			setter = func(target reflect.Value, value string) error {
				v, err := parse(value)
				if err != nil {
					return err
				}
				if s, isString := v.(string); isString {
					// Go through Set() or UnmarshalText()
					return setString(target, s)
				}
				//fmt.Printf("%T -> %s\n", v, target.Type())
				//fmt.Printf("%s -> %s\n", reflect.ValueOf(&v).Elem().Elem().Type(), target.Type())
				target.Set(reflect.ValueOf(&v).Elem().Elem())
				return nil
			}
		}
	}
	return &slice{v, itemType, separator, setter}
}

type slice struct {
	slice     reflect.Value
	itemType  reflect.Type
	separator string
	setString func(reflect.Value, string) error
}

func (sl *slice) String() string {
	return ""
}

func (sl *slice) appnd(s string) error {
	v := reflect.New(sl.itemType)
	err := sl.setString(v.Elem(), s)
	if err != nil {
		return err
	}
	sl.slice.Set(reflect.Append(sl.slice, v.Elem()))
	return nil
}

func (sl *slice) Set(s string) error {
	if len(sl.separator) != 0 {
		for _, item := range strings.Split(s, sl.separator) {
			if err := sl.appnd(item); err != nil {
				return err
			}
		}
		return nil
	}
	return sl.appnd(s)
}

func (sl *slice) Get() interface{} {
	return sl.slice.Interface()
}
