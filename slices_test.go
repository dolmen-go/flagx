package flagx_test

import (
	"flag"
	"strconv"
	"testing"

	"github.com/dolmen-go/flagx"
)

func checkIntSlice(tester *varTester) {
	tester.CheckParse([]string{}, ([]int)(nil))
	tester.CheckParse([]string{"a"}, ([]int)(nil))
	tester.CheckParse([]string{"-ints", "0"}, []int{0})
	tester.CheckParse([]string{"-ints", "1,2"}, []int{1, 2})
	tester.CheckParse([]string{"-ints", "2", "-ints", "3"}, []int{2, 3})
	tester.CheckParse([]string{"-ints", "1,2,3"}, []int{1, 2, 3})
	tester.CheckParse([]string{"-ints", "1,2,3", "-ints", "4"}, []int{1, 2, 3, 4})
	tester.CheckParse([]string{"-ints", "1,2,3", "-ints", "4,5"}, []int{1, 2, 3, 4, 5})
	tester.CheckParse([]string{"-ints", "1,2,3", "-ints", "4,5,6"}, []int{1, 2, 3, 4, 5, 6})
	tester.CheckParse([]string{"-ints", "0xf,010,-1"}, []int{15, 8, -1})
	tester.CheckParse([]string{"-ints", "0x7fffffff"}, []int{0x7fffffff})
	tester.CheckParse([]string{"-ints", "-0x80000000"}, []int{-0x80000000})

	tester.CheckHelp()
}

func checkStringSlice(tester *varTester) {
	tester.CheckParse([]string{}, ([]string)(nil))
	tester.CheckParse([]string{"a"}, ([]string)(nil))
	tester.CheckParse([]string{"-strings", "a"}, []string{"a"})
	tester.CheckParse([]string{"-strings", "a,b"}, []string{"a", "b"})
	tester.CheckParse([]string{"-strings", "a", "-strings", "b"}, []string{"a", "b"})
}

func TestIntSlice(t *testing.T) {
	checkIntSlice(&varTester{
		t:        t,
		flagName: "ints",
		buildVar: func() (flag.Getter, interface{}) {
			var value []int
			return flagx.IntSlice{&value}, &value
		}})
}

func TestSlice(t *testing.T) {
	checkIntSlice(&varTester{
		t:        t,
		flagName: "ints",
		buildVar: func() (flag.Getter, interface{}) {
			var value []int
			return flagx.Slice(&value, ",", func(s string) (interface{}, error) {
				n, err := strconv.ParseInt(s, 0, 0)
				if err != nil {
					return nil, nil
				}
				return int(n), nil
			}), &value
		}})
	checkStringSlice(&varTester{
		t:        t,
		flagName: "strings",
		buildVar: func() (flag.Getter, interface{}) {
			var value []string
			return flagx.Slice(&value, ",", func(s string) (interface{}, error) {
				return s, nil
			}), &value
		}})
	checkStringSlice(&varTester{
		t:        t,
		flagName: "strings",
		buildVar: func() (flag.Getter, interface{}) {
			var value []string
			return flagx.Slice(&value, ",", nil), &value
		}})
}
