package flagx_test

import (
	"flag"
	"fmt"

	"github.com/dolmen-go/flagx"
)

func ExampleYesNo() {
	flags := flag.NewFlagSet("test", flag.PanicOnError) // Usually flag.CommandLine

	b1 := false
	b2 := true
	flags.Var(flagx.YesNo(&b1), "opt1", "boolean option")
	flags.Var(flagx.YesNo(&b2), "opt2", "boolean option")

	flags.Parse([]string{"-opt1=yes", "-opt2=no"})

	fmt.Println(b1)
	fmt.Println(b2)

	// Output:
	// true
	// false
}
