package flagx_test

import (
	"flag"
	"testing"

	"go.uber.org/thriftrw/ptr"

	"github.com/dolmen-go/flagx"
)

func TestNullString(t *testing.T) {
	tester := varTester{
		t:        t,
		flagName: "value",
		buildVar: func() (flag.Getter, interface{}) {
			var value *string
			return flagx.NullString(&value), &value
		}}

	tester.Check([]string{}, (*string)(nil))
	tester.Check([]string{"a"}, (*string)(nil))
	tester.Check([]string{"-value", "x"}, ptr.String("x"))
	tester.Check([]string{"-value", ""}, ptr.String(""))
	tester.Check([]string{"-value", "a", "-value", "b"}, ptr.String("b"))
}
