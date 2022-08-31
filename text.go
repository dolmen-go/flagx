package flagx

type textValue struct {
	value interface {
		// encoding.TextMarshaler
		MarshalText() (text []byte, err error)
		// encoding.TextUnmarshaler
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

// Text wraps a Text{Unm,M}arshaler as a flag.Getter
// which can then be passed to flag.Var() / flag.FlagSet.Var()
//
// Deprecated: use flag.TextVar instead available since Go 1.19.
func Text(v interface {
	// encoding.TextMarshaler
	MarshalText() (text []byte, err error)
	// encoding.TextUnmarshaler
	UnmarshalText(text []byte) error
}) Value {
	return &textValue{v}
}
