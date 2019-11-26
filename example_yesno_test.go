package flagx_test

import (
	"flag"
	"fmt"

	"github.com/dolmen-go/flagx"
)

func ExampleYesNo() {
	flags := flag.FlagSet{} // Usually flags.CommandLine

	b1 := false
	b2 := true
	flags.Var(flagx.YesNo(&b1), "opt1", "boolean option")
	flags.Var(flagx.YesNo(&b2), "opt2", "boolean option")

	if err := flags.Parse([]string{"-opt1=yes", "-opt2=no"}); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(b1)
	fmt.Println(b2)

	// Output:
	// true
	// false
}
