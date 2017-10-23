package flagx_test

import (
	"bytes"
	"flag"
	"reflect"
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
	output   bytes.Buffer
}

func (tester *varTester) Check(args []string, expected interface{}) {
	flagValue, pvalue := tester.buildVar()
	if reflect.TypeOf(pvalue).Kind() != reflect.Ptr {
		panic("varBuilder must return a pointer")
	}

	flags := flag.NewFlagSet("Test", flag.ContinueOnError)
	tester.output.Reset()
	flags.SetOutput(&tester.output)
	if args == nil {
		args = []string{}
	}

	flags.Var(flagValue, tester.flagName, "Value")
	if err := flags.Parse(args); err != nil {
		tester.t.Fatalf("Unexpected error: %s\nError output:\n%s", err, tester.output.String())
	}
	// Dereference pvalue
	value := reflect.ValueOf(pvalue).Elem().Interface()
	if !reflect.DeepEqual(value, expected) {
		tester.t.Errorf("got %#v expected %#v", value, expected)
	}
	if tester.output.Len() > 0 {
		tester.t.Errorf("Error output:\n%s", tester.output.String())
	}
}
