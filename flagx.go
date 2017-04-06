package flagtext

// Value is like flag.Getter
type Value interface {
	String() string
	Set(string) error
	Get() interface{}
}
