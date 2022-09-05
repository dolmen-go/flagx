package flagnet

// Value is like [flag.Getter] (which is a superset of [flag.Value]).
type Value interface {
	String() string
	Set(string) error
	Get() interface{}
}
