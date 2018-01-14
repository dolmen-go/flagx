package flagfile

import (
	"flag"
	"fmt"
	"io"
	"os"
	"reflect"
)

// Decoder can decode a stream of bytes into structured values.
// The root value for flagfile.File must be of type:
//     []string
//     []interface{}
//     map[string]interface{}
//
// Example: encoding/json.Decoder
type Decoder interface {
	Decode(interface{}) error
}

type DecoderBuilder func(r io.Reader) Decoder

type file struct {
	flagset *flag.FlagSet
	decoder DecoderBuilder
}

func File(flagset *flag.FlagSet, decoder DecoderBuilder) flag.Value {
	return &file{
		flagset: flagset,
		decoder: decoder,
	}
}

func (file) String() string { return "" }

func (f *file) Set(filepath string) error {
	r, err := os.Open(filepath)
	if err != nil {
		// TODO return a structured error to allow l18n
		return fmt.Errorf("%s: %s", filepath, err)
	}
	dec := f.decoder(r)
	var v interface{}
	err = dec.Decode(&v)
	if err != nil {
		// TODO return a structured error to allow l18n
		return fmt.Errorf("%s: can't decode: %s", filepath, err)
	}

	switch v := v.(type) {
	case []interface{}:
		args := make([]string, 0, len(v))
		for _, arg := range v {
			args = append(args, fmt.Sprint(arg))
		}
		err = f.flagset.Parse(args)
	case []string:
		err = f.flagset.Parse(v)
	case map[string]interface{}:
		err = parseObject(f.flagset, v)
		// TODO map[string]string
		// TODO map[string][]string
	default:
		return fmt.Errorf("unexpected type %T", v)
	}
	if err != nil {
		err = fmt.Errorf("%s#%s", filepath, err)
	}
	return err
}

func parseObject(flagset *flag.FlagSet, m map[string]interface{}) error {
	for k, v := range m {
		f := flagset.Lookup(k)
		if f == nil {
			return fmt.Errorf("/%s: unknown flag", k)
		}
		var err error
		switch v := v.(type) {
		case string:
			err = f.Value.Set(v)
		case []string:
			for i, v := range v {
				err = f.Value.Set(v)
				if err != nil {
					return fmt.Errorf("/%s/%d: %s", k, i, err)
				}
			}
		default:
			rv := reflect.ValueOf(v)
			if rv.Kind() == reflect.Array || rv.Kind() == reflect.Slice {
				for i := 0; i < rv.Len(); i++ {
					err = f.Value.Set(fmt.Sprint(rv.Index(i).Interface()))
					if err != nil {
						return fmt.Errorf("/%s/%d: %s", k, i, err)
					}
				}
			} else {
				err = f.Value.Set(fmt.Sprint(v))
			}
		}
		if err != nil {
			return fmt.Errorf("/%s: %s", k, err)
		}
	}

	return nil
}
