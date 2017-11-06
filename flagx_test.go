package flagx_test

import (
	"bytes"
	"flag"
	"reflect"
	"strings"
	"testing"

	"github.com/dolmen-go/flagx"
)

// Check that our flagx.Value is the same as flag.Getter
var _ flag.Getter = flagx.Value(nil)

type varBuilder func() (flag.Getter, interface{})

type varTester struct {
	t        *testing.T
	flagName string
	buildVar varBuilder
}

func (tester *varTester) CheckParse(args []string, expected interface{}) {
	flagValue, pvalue := tester.buildVar()
	kind := reflect.TypeOf(pvalue).Kind()
	if kind != reflect.Ptr && kind != reflect.Map {
		panic("varBuilder must return a pointer")
	}

	flags := flag.NewFlagSet("Test", flag.ContinueOnError)
	var output bytes.Buffer
	flags.SetOutput(&output)
	if args == nil {
		args = []string{}
	}

	flags.Var(flagValue, tester.flagName, "Value")

	if err := flags.Parse(args); err != nil {
		tester.t.Fatalf("Unexpected error: %s\nError output:\n%s", err, output.String())
	}
	// Dereference pvalue
	var value interface{}
	switch kind {
	case reflect.Ptr:
		value = reflect.ValueOf(pvalue).Elem().Interface()
	case reflect.Map:
		value = pvalue
	}
	if !reflect.DeepEqual(value, expected) {
		tester.t.Errorf("got %#v expected %#v", value, expected)
	}
	flgVar := flags.Lookup(tester.flagName).Value.(flag.Getter)
	valueFromFlag := flgVar.Get()
	if !reflect.DeepEqual(valueFromFlag, expected) {
		tester.t.Errorf("got %#v expected %#v", valueFromFlag, expected)
	}
	_ = flgVar.String()

	if output.Len() > 0 {
		tester.t.Errorf("Error output:\n%s", output.String())
	}
}

func (tester *varTester) CheckHelp() {
	flagValue, pvalue := tester.buildVar()
	if reflect.TypeOf(pvalue).Kind() != reflect.Ptr {
		panic("varBuilder must return a pointer")
	}

	flags := flag.NewFlagSet("Test", flag.ContinueOnError)
	var output bytes.Buffer
	flags.SetOutput(&output)

	flags.Var(flagValue, tester.flagName, "set arg `v`")

	err := flags.Parse([]string{"-h"})
	if err != flag.ErrHelp {
		tester.t.Fatalf("ErrHelp expected, got %q", err)
	}

	out := output.String()
	if !strings.Contains(out, "-"+tester.flagName) {
		tester.t.Errorf("Incorrect usage message: expected mention of `-%s`, but got:\n%s", tester.flagName, out)
	} else {
		tester.t.Logf("Help message:\n%s", out)
	}
}
