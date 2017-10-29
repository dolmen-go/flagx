package flagx_test

import (
	"flag"
	"testing"

	"github.com/dolmen-go/flagx"
)

func TestIntSlice(t *testing.T) {
	tester := varTester{
		t:        t,
		flagName: "ints",
		buildVar: func() (flag.Getter, interface{}) {
			var value []int
			return flagx.IntSlice{&value}, &value
		}}

	tester.CheckParse([]string{}, ([]int)(nil))
	tester.CheckParse([]string{"a"}, ([]int)(nil))
	tester.CheckParse([]string{"-ints", "0"}, []int{0})
	tester.CheckParse([]string{"-ints", "1,2"}, []int{1, 2})
	tester.CheckParse([]string{"-ints", "2", "-ints", "3"}, []int{2, 3})
	tester.CheckParse([]string{"-ints", "1,2,3"}, []int{1, 2, 3})
	tester.CheckParse([]string{"-ints", "1,2,3", "-ints", "4"}, []int{1, 2, 3, 4})
	tester.CheckParse([]string{"-ints", "1,2,3", "-ints", "4,5"}, []int{1, 2, 3, 4, 5})
	tester.CheckParse([]string{"-ints", "1,2,3", "-ints", "4,5,6"}, []int{1, 2, 3, 4, 5, 6})
	tester.CheckParse([]string{"-ints", "0xf"}, []int{15})

	tester.CheckHelp()
}
