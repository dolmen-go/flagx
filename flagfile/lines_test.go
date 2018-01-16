package flagfile_test

import (
	"flag"
	"fmt"

	"github.com/dolmen-go/flagx/flagfile"
)

var _ flagfile.Loader = flagfile.Lines

func saveCommandLine() func() {
	backup := flag.CommandLine
	flag.CommandLine = flag.NewFlagSet("", flag.ContinueOnError)
	return func() {
		flag.CommandLine = backup
	}
}

func Example() {
	defer saveCommandLine()() // Only for the example context

	flag.Var(flagfile.File(flag.CommandLine, flagfile.Lines), "args", "args from file (one line = one argument)")
	b := flag.Bool("bool", false, "bool test")
	i := flag.Int("int", 0, "int test")
	s := flag.String("string", "", "string test")

	err := flag.CommandLine.Parse([]string{
		"-args", "testdata/args.txt",
	})
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(*b, *i, *s)

	// Output:
	// true 1 foo
}
