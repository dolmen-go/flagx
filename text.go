package flagx

type textValue struct {
	value interface {
		// See [encoding.TextMarshaler].
		MarshalText() (text []byte, err error)
		// See [encoding.TextUnmarshaler].
		UnmarshalText(text []byte) error
	}
}

func (v *textValue) String() string {
	b, err := v.value.MarshalText()
	if err != nil {
		// Panic?
		return ""
	}
	return string(b)
}

func (v *textValue) Set(str string) error {
	return v.value.UnmarshalText([]byte(str))
}

func (v *textValue) Get() interface{} {
	return v.value
}

// Text wraps an [encoding.TextUnmarshaler] + [encoding.TextMarshaler] as a [flag.Getter]
// which can then be passed to [flag.Var] / [flag.FlagSet.Var].
//
// Note: you might prefer to use [flag.TextVar] which is available since Go 1.19.
func Text(v interface {
	// See [encoding.TextMarshaler].
	MarshalText() (text []byte, err error)
	// See [encoding.TextUnmarshaler].
	UnmarshalText(text []byte) error
}) Value {
	return &textValue{v}
}
