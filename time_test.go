package flagx_test

import (
	"flag"
	"testing"
	"time"

	"github.com/dolmen-go/flagx"
)

func TestTime(t *testing.T) {
	tester := varTester{
		t:        t,
		flagName: "when",
		buildVar: func() (flag.Getter, interface{}) {
			var t time.Time
			return flagx.Time{&t, "2006-01-02T15:04:05", time.Local}, &t
		}}

	tester.CheckParse([]string{}, time.Time{})
	tester.CheckParse([]string{"a"}, time.Time{})
	tester.CheckParse([]string{"-when", "2017-10-29T14:31:22"}, time.Date(2017, 10, 29, 14, 31, 22, 0, time.Local))

	tester.CheckHelp()
}
