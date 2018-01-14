package flagjson

import (
	"encoding/json"
	"flag"
	"io"
	"reflect"

	"github.com/dolmen-go/flagx/flagfile"
	"github.com/dolmen-go/jsonptr"
)

// File allows to declare a command-line flag that will inject
// the JSON content of the filename given by the user as arguments
// of the command-line.
//
// The pointer argument allows to load the arguments from a part of the file by
// giving a JSON Pointer (RFC 6901) of the location. "" means the root
// value, which is the whole file.
func File(flagset *flag.FlagSet, pointer string) flag.Value {
	var decoder flagfile.DecoderBuilder
	if pointer == "" {
		decoder = func(r io.Reader) flagfile.Decoder {
			return json.NewDecoder(r)
		}
	} else {
		ptr, err := jsonptr.Parse(pointer)
		if err != nil {
			panic(err)
		}
		decoder = func(r io.Reader) flagfile.Decoder {
			return &partialDecoder{
				decoder: json.NewDecoder(r),
				pointer: ptr,
			}
		}
	}
	return flagfile.File(flagset, decoder)
}

type partialDecoder struct {
	decoder flagfile.Decoder
	pointer jsonptr.Pointer
}

func (d *partialDecoder) Decode(v interface{}) error {
	value, err := d.pointer.In(d.decoder)
	if err != nil {
		return err
	}
	reflect.ValueOf(v).Elem().Set(reflect.ValueOf(value))
	return nil
}
