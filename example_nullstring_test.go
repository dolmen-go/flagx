package flagx_test

import (
	"flag"
	"fmt"

	"github.com/dolmen-go/flagx"
)

func ExampleNullString() {
	flags := flag.NewFlagSet("test", flag.PanicOnError) // Usually flag.CommandLine

	var value *string
	flags.Var(flagx.NullString(&value), "value", "value")

	flags.Parse([]string{"-value", "hello"})

	if value == nil {
		fmt.Println("not initialised")
	} else {
		fmt.Printf("%q\n", *value)
	}

	// Output:
	// "hello"
}
