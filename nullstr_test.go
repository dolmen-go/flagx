package flagx_test

import (
	"flag"
	"reflect"
	"testing"

	"go.uber.org/thriftrw/ptr"

	"github.com/dolmen-go/flagx"
)

type nullStringTester struct {
	t        *testing.T
	flagName string
}

func (nst *nullStringTester) Check(args []string, expected *string) {
	var value *string
	flags := flag.NewFlagSet("TestNullString", flag.ContinueOnError)
	if args == nil {
		args = []string{}
	}
	flags.Var(flagx.NullString(&value), nst.flagName, "Value")
	if err := flags.Parse(args); err != nil {
		nst.t.Fatalf("Unexpected error: %s", err)
	}
	if !reflect.DeepEqual(value, expected) {
		nst.t.Errorf("got %#v expected %#v", value, expected)
	}
}

func TestNullString(t *testing.T) {
	tns := nullStringTester{t, "value"}
	tns.Check([]string{}, nil)
	tns.Check([]string{"a"}, nil)
	tns.Check([]string{"-value", "x"}, ptr.String("x"))
	tns.Check([]string{"-value", ""}, ptr.String(""))
	tns.Check([]string{"-value", "a", "-value", "b"}, ptr.String("b"))
}
