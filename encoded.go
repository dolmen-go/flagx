package flagx

// Encoding is an interface implemented by stdlib encoding/base{32,64}.Encoding
type Encoding interface {
	DecodeString(string) ([]byte, error)
	EncodeToString(src []byte) string
}

type encodedValue struct {
	encoding Encoding
	value    *[]byte
}

func (v *encodedValue) String() string {
	return v.encoding.EncodeToString(*v.value)
}

func (v *encodedValue) Set(str string) (err error) {
	*v.value, err = v.encoding.DecodeString(str)
	return
}

func (v *encodedValue) Get() interface{} {
	return *v.value
}

// Encoded wraps a reference to a []byte as a flag.Value
func Encoded(value *[]byte, encoding Encoding) Value {
	return &encodedValue{value: value, encoding: encoding}
}
