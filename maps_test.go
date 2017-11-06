package flagx_test

import (
	"flag"
	"strconv"
	"testing"

	"github.com/dolmen-go/flagx"
)

func checkIntMap(tester *varTester) {
	tester.CheckParse([]string{}, map[string]int{})
	tester.CheckParse([]string{"-kv", "a=0"}, map[string]int{"a": 0})
	tester.CheckParse([]string{"-kv", "a=0", "-kv", "b=1"}, map[string]int{"a": 0, "b": 1})
}

func TestMap(t *testing.T) {
	checkIntMap(&varTester{
		t:        t,
		flagName: "kv",
		buildVar: func() (flag.Getter, interface{}) {
			m := make(map[string]int)
			return flagx.Map(m, func(s string) (interface{}, error) {
				n, err := strconv.ParseInt(s, 0, 0)
				if err != nil {
					return nil, nil
				}
				return int(n), nil
			}), m
		}})
}
