package flagx

import (
	"strconv"
)

// Func wraps a function like [flag.Value.Set] as a [flag.Value].
type Func func(s string) error

func (f Func) Set(s string) error { return f(s) }

func (Func) String() string { return "" }

func (Func) Get() interface{} { return nil }

// BoolFunc wraps a function as a boolean [flag.Value].
type BoolFunc func(b bool) error

func (BoolFunc) IsBoolFlag() bool { return true }

func (f BoolFunc) Set(s string) error {
	v, err := strconv.ParseBool(s)
	if err != nil {
		return err
	}
	return f(v)
}

func (BoolFunc) String() string { return "" }

func (BoolFunc) Get() interface{} { return nil }
