package flagtext

type varText struct {
	value interface {
		// encoding.TextMarshaler
		MarshalText() (text []byte, err error)
		// encoding.TextUnmarshaler
		UnmarshalText(text []byte) error
	}
}

func (v *varText) String() string {
	b, err := v.value.MarshalText()
	if err != nil {
		// Panic?
		return ""
	}
	return string(b)
}

func (v *varText) Set(str string) error {
	return v.value.UnmarshalText([]byte(str))
}

func (v *varText) Get() interface{} {
	return v.value
}

// VarText wraps a Text{Unm,M}arshaler as a flag.Getter
// which can then be passed to flag.Var() / flag.FlagSet.Var()
func VarText(v interface {
	// encoding.TextMarshaler
	MarshalText() (text []byte, err error)
	// encoding.TextUnmarshaler
	UnmarshalText(text []byte) error
}) interface {
	String() string
	Set(string) error
	Get() interface{}
} {
	return &varText{v}
}
