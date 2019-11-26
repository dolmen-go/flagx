package flagfile_test

import (
	"flag"
	"fmt"

	"github.com/dolmen-go/flagx/flagfile"
)

func Example() {
	var cmdLine flag.FlagSet // Usually flags.CommandLine

	cmdLine.Var(flagfile.File(&cmdLine, flagfile.Lines), "args", "args from file (one line = one argument)")
	b := cmdLine.Bool("bool", false, "bool test")
	i := cmdLine.Int("int", 0, "int test")
	s := cmdLine.String("string", "", "string test")

	err := cmdLine.Parse([]string{
		"-args", "testdata/args.txt",
	})
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(*b, *i, *s)

	// Output:
	// true 1 foo
}
