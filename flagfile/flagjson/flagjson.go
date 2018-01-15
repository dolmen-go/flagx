package flagjson

import (
	"encoding/json"
	"flag"
	"io"

	"github.com/dolmen-go/flagx/flagfile"
	"github.com/dolmen-go/jsonptr"
)

// File allows to declare a command-line flag that will inject
// the JSON content of the filename given by the user as arguments
// of the command-line.
//
// The filename may include an URI fragment starting with '#' and followed
// by a JSON Pointer (RFC 6901) to load the arguments from a part of the file.
// The empty fragment ("") means the root value, which is the whole file.
func File(flagset *flag.FlagSet) flag.Value {
	return flagfile.File(flagset, loader)
}

func loader(r io.Reader, fragment string) (interface{}, error) {
	decoder := json.NewDecoder(r)

	if fragment == "" {
		var v interface{}
		return v, decoder.Decode(&v)
	}

	return jsonptr.Get(decoder, fragment)
}
