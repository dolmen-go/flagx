package flagx_test

import (
	"bytes"
	"flag"
	"reflect"
	"testing"

	"go.uber.org/thriftrw/ptr"

	"github.com/dolmen-go/flagx"
)

type nullStringTester struct {
	t        *testing.T
	flagName string
	output   bytes.Buffer
}

func (nst *nullStringTester) Check(args []string, expected *string) {
	var value *string
	flags := flag.NewFlagSet("TestNullString", flag.ContinueOnError)
	flags.SetOutput(&nst.output)
	if args == nil {
		args = []string{}
	}
	flags.Var(flagx.NullString(&value), nst.flagName, "Value")
	if err := flags.Parse(args); err != nil {
		nst.t.Fatalf("Unexpected error: %s\nError output:\n%s", err, nst.output.String())
	}
	if !reflect.DeepEqual(value, expected) {
		nst.t.Errorf("got %#v expected %#v", value, expected)
	}
	if nst.output.Len() > 0 {
		nst.t.Errorf("Error output:\n%s", nst.output.String())
		nst.output.Reset()
	}
}

func TestNullString(t *testing.T) {
	tns := nullStringTester{t: t, flagName: "value"}
	tns.Check([]string{}, nil)
	tns.Check([]string{"a"}, nil)
	tns.Check([]string{"-value", "x"}, ptr.String("x"))
	tns.Check([]string{"-value", ""}, ptr.String(""))
	tns.Check([]string{"-value", "a", "-value", "b"}, ptr.String("b"))
}
