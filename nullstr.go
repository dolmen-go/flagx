package flagx

// NullString allow to define string flags when you want to distinguish
// the case where the flag is never set from when it is set to an empty string
func NullString(value **string) Value {
	return nullStr{value}
}

type nullStr struct {
	Pointer **string
}

func (ns nullStr) String() string {
	if ns.Pointer == nil {
		// When called by flag.isZeroValue
		return ""
	}
	if *ns.Pointer == nil {
		return ""
	}
	return **ns.Pointer
}

func (ns nullStr) Set(str string) (err error) {
	if *ns.Pointer == nil {
		*ns.Pointer = new(string)
	}
	**ns.Pointer = str
	return nil
}

func (ns nullStr) Get() interface{} {
	return *ns.Pointer
}
