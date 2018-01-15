package flagfile

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
)

// FlagSet is a subset of methods of *flag.FlagSet
type FlagSet interface {
	Parse([]string) error
	Lookup(name string) *flag.Flag
}

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
	flagset    FlagSet
	decoder    DecoderBuilder
	contextDir []string
}

// File allows to define a command-line flag that gives a path to a structured
// file whose content will be expanded as command-line arguments and injected into flagset.
//
// The structured data may be:
// - an array of strings given to flagset.Parse
// - a map where keys are argument names (without leading '-') and keys are values
//   or arrays of values
//
// This flag is reentrant, so the file may refer to (include) other files by reusing the
// same flag.
func File(flagset FlagSet, decoder DecoderBuilder) flag.Value {
	return &file{
		flagset: flagset,
		decoder: decoder,
	}
}

func (file) String() string { return "" }

// Set loads the given file, decodes it and interprets
// its structured data as a list of arguments.
//
// Set is part of the flag.Value interface.
func (f *file) Set(path string) error {
	if !filepath.IsAbs(path) {
		var ctxDir string
		if len(f.contextDir) == 0 {
			var e error
			ctxDir, e = os.Getwd()
			if e != nil {
				return e
			}
		} else {
			ctxDir = f.contextDir[len(f.contextDir)-1]
		}
		path = filepath.Join(ctxDir, path)
	}

	r, err := os.Open(path)
	if err != nil {
		// TODO return a structured error to allow l18n
		return fmt.Errorf("%s: %s", path, err)
	}

	v, err := func(r *os.File) (interface{}, error) {
		defer r.Close()
		dec := f.decoder(r)
		var v interface{}
		err = dec.Decode(&v)
		return v, err
	}(r)
	if err != nil {
		// TODO return a structured error to allow l18n
		return fmt.Errorf("%s: can't decode: %s", path, err)
	}

	// push directory of the path
	f.contextDir = append(f.contextDir, filepath.Dir(path))
	defer func() {
		// pop contextDir
		f.contextDir = f.contextDir[:len(f.contextDir)-1]
	}()

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
		err = fmt.Errorf("%s#%s", path, err)
	}
	return err
}

func parseObject(flagset FlagSet, m map[string]interface{}) error {
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
