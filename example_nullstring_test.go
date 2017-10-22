package flagx_test

import (
	"flag"
	"fmt"
	"os"

	"github.com/dolmen-go/flagx"
)

func ExampleNullString() {
	// Overwrite os.Args just for testing
	savedArgs := os.Args
	defer func() { os.Args = savedArgs }()
	os.Args = []string{"example", "-value", "hello"}

	var value *string
	flag.Var(flagx.NullString(&value), "value", "value")

	flag.Parse()

	if value == nil {
		fmt.Println("not initialised")
	} else {
		fmt.Printf("%q\n", *value)
	}

	// Output:
	// "hello"
}
