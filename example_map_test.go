package flagx_test

import (
	"flag"
	"fmt"
	"strconv"

	"github.com/dolmen-go/flagx"
)

func ExampleMap() {
	flags := flag.NewFlagSet("test", flag.PanicOnError) // Usually flag.CommandLine

	m := make(map[string]int)
	flags.Var(flagx.Map(m, func(s string) (interface{}, error) {
		return strconv.Atoi(s)
	}), "define", "define key=value pairs")

	flags.Parse([]string{"-define", "x=4", "-define", "y=5"})

	if len(m) == 0 {
		fmt.Println("not initialised")
	} else {
		fmt.Printf("x=%d, y=%d\n", m["x"], m["y"])
	}

	// Output:
	// x=4, y=5
}
