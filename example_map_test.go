package flagx_test

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/dolmen-go/flagx"
)

func ExampleMap() {
	// Overwrite os.Args just for testing
	savedArgs := os.Args
	defer func() { os.Args = savedArgs }()
	os.Args = []string{"example", "-define", "x=4", "-define", "y=5"}

	m := make(map[string]int)
	flag.Var(flagx.Map(m, func(s string) (interface{}, error) {
		return strconv.Atoi(s)
	}), "define", "define key=value pairs")

	flag.Parse()

	if len(m) == 0 {
		fmt.Println("not initialised")
	} else {
		fmt.Printf("x=%d, y=%d\n", m["x"], m["y"])
	}

	// Output:
	// x=4, y=5
}
