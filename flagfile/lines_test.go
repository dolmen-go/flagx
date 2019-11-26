package flagfile_test

import (
	"flag"
	"fmt"

	"github.com/dolmen-go/flagx/flagfile"
)

func Example() {
	flags := flag.NewFlagSet("test", flag.PanicOnError) // Usually flag.CommandLine

	flags.Var(flagfile.File(flags, flagfile.Lines), "args", "args from file (one line = one argument)")
	b := flags.Bool("bool", false, "bool test")
	i := flags.Int("int", 0, "int test")
	s := flags.String("string", "", "string test")

	flags.Parse([]string{
		"-args", "testdata/args.txt",
	})

	fmt.Println(*b, *i, *s)

	// Output:
	// true 1 foo
}
