package flagx

import "os"

// Env wraps a [flag.Value] and initializes it with the value
// of the given environment variable if set.
// The lookup occurs immediately so it happens before command line parsing.
func Env(key string, v Value) Value {
	if x, found := os.LookupEnv(key); found {
		// Ignore errors
		_ = v.Set(x)
	}
	return v
}
