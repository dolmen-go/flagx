package flagx_test

import (
	"flag"
	"fmt"
	"time"

	"github.com/dolmen-go/flagx"
)

func ExampleTime() {
	flags := flag.NewFlagSet("test", flag.PanicOnError) // Usually flag.CommandLine

	var t time.Time

	flags.Var(flagx.Time{
		Time:     &t,
		Format:   time.RFC3339,
		Location: time.UTC,
	}, "time", "time value")

	flags.Parse([]string{"-time=2006-01-02T15:04:05Z"})

	fmt.Println(t)

	// Output:
	// 2006-01-02 15:04:05 +0000 UTC
}
