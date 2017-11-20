package flagtrace_test

import (
	"flag"
	"fmt"

	"github.com/dolmen-go/flagx/flagtrace"
)

func ExampleRegister() {
	stopTracing := flagtrace.Register(flag.CommandLine, "debug.trace", "trace `file` (for go tool trace)")
	defer stopTracing()

	flag.Parse()

	fmt.Println("hello, world!")

	// Output:
	// hello, world!
}
