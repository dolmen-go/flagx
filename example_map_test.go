package flagx_test

import (
	"flag"
	"fmt"
	"strconv"

	"github.com/dolmen-go/flagx"
)

func ExampleMap() {
	flags := flag.FlagSet{} // Usually flag.CommandLine

	m := make(map[string]int)
	flags.Var(flagx.Map(m, func(s string) (interface{}, error) {
		return strconv.Atoi(s)
	}), "define", "define key=value pairs")

	if err := flags.Parse([]string{"-define", "x=4", "-define", "y=5"}); err != nil {
		fmt.Println(err)
	}

	if len(m) == 0 {
		fmt.Println("not initialised")
	} else {
		fmt.Printf("x=%d, y=%d\n", m["x"], m["y"])
	}

	// Output:
	// x=4, y=5
}
