package flagx_test

import (
	"flag"
	"testing"

	"github.com/dolmen-go/flagx"
)

func stringPtr(s string) *string {
	return &s
}

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
	tester.CheckParse([]string{"-value", "x"}, stringPtr("x"))
	tester.CheckParse([]string{"-value", ""}, stringPtr(""))
	tester.CheckParse([]string{"-value", "a", "-value", "b"}, stringPtr("b"))

	tester.CheckHelp()
}
