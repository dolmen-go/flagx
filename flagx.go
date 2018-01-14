package flagx

// Value is like flag.Getter (which is a superset of flag.Value)
type Value interface {
	String() string
	Set(string) error
	Get() interface{}
}

// Dummy is a flag.Value that does nothing.
type Dummy struct{}

func (Dummy) String() string { return "" }

func (Dummy) Set(s string) error { return nil }

func (Dummy) Get() interface{} { return nil }
