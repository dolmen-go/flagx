package flagfile

import (
	"flag"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

// FlagSet is a subset of methods of *flag.FlagSet
type FlagSet interface {
	Parse([]string) error
	Args() []string
	Set(name, value string) error
}

// Loader loads structured data from r.
// The fragment may point to a partial of the content.
//
// The output is expected to be:
// - a slice of strings
// - a map where keys are strings and values are either string or slice of strings
type Loader func(r io.Reader, fragment string) (interface{}, error)

type file struct {
	flagset    FlagSet
	load       Loader
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
func File(flagset FlagSet, loader Loader) flag.Value {
	return &file{
		flagset: flagset,
		load:    loader,
	}
}

func (file) String() string { return "" }

// Set loads the given file, decodes it and interprets
// its structured data as a list of arguments.
//
// Set is part of the flag.Value interface.
func (f *file) Set(path string) error {
	var fragment string
	if i := strings.IndexByte(path, '#'); i > 0 {
		var err error
		if fragment, err = url.QueryUnescape(path[i+1:]); err != nil {
			return fmt.Errorf("invalid fragment %q", path[i+1:])
		}
		path = path[:i]
	}

	// Resolve the path to an absolute path using either the path of the enclosing file
	// or the current working directory.
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

	v, err := func(r *os.File, load Loader) (interface{}, error) {
		defer r.Close()
		return load(r, fragment)
	}(r, f.load)
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
	case nil: // NOOP
	case []interface{}:
		if len(v) == 0 {
			break
		}
		args := make([]string, 0, len(v))
		for _, arg := range v {
			args = append(args, fmt.Sprint(arg))
		}
		remaining := f.flagset.Args()
		err = f.flagset.Parse(args)
		if err != nil {
			return err
		}
		err = f.flagset.Parse(remaining)
	case []string:
		if len(v) > 0 {
			remaining := f.flagset.Args()
			err = f.flagset.Parse(v)
			if err != nil {
				return err
			}
			err = f.flagset.Parse(remaining)
		}
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
		var err error
		switch v := v.(type) {
		case string:
			err = flagset.Set(k, v)
		case []string:
			for i, v := range v {
				err = flagset.Set(k, v)
				if err != nil {
					return fmt.Errorf("/%s/%d: %s", k, i, err)
				}
			}
		default:
			rv := reflect.ValueOf(v)
			if rv.Kind() == reflect.Array || rv.Kind() == reflect.Slice {
				for i := 0; i < rv.Len(); i++ {
					err = flagset.Set(k, fmt.Sprint(rv.Index(i).Interface()))
					if err != nil {
						return fmt.Errorf("/%s/%d: %s", k, i, err)
					}
				}
			} else {
				err = flagset.Set(k, fmt.Sprint(v))
			}
		}
		if err != nil {
			return fmt.Errorf("/%s: %s", k, err)
		}
	}

	return nil
}
