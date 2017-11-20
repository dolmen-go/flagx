package main

import (
	"flag"
	"time"

	"github.com/dolmen-go/flagx/flagtrace"
)

func main() {
	stopTracing := flagtrace.Register(flag.CommandLine, "debug.trace", "trace `file` (for go tool trace)")
	defer stopTracing()

	flag.Parse()

	time.Sleep(1 * time.Second)
}
