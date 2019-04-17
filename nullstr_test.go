package flagx_test

import (
	"flag"
	"testing"

	"github.com/wacul/ptr"

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

	tester.CheckParse([]string{}, (*string)(nil))
	tester.CheckParse([]string{"a"}, (*string)(nil))
	tester.CheckParse([]string{"-value", "x"}, ptr.String("x"))
	tester.CheckParse([]string{"-value", ""}, ptr.String(""))
	tester.CheckParse([]string{"-value", "a", "-value", "b"}, ptr.String("b"))

	tester.CheckHelp()
}
