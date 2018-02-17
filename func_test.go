package flagx_test

import (
	"encoding/hex"
	"flag"
	"fmt"
	"math/rand"
	"os"

	"github.com/dolmen-go/flagx"
)

/*
	defer func(args []string, cmdline *flag.FlagSet) {
		os.Args = args
		flag.CommandLine = cmdline
	}(os.Args, flag.CommandLine)
*/

func ExampleFunc() {
	os.Args = []string{os.Args[0], "-push=a", "-push=b"}
	flag.CommandLine = flag.NewFlagSet("test", flag.PanicOnError)

	var all []string

	// A flag that apppends the given value to slice all
	flag.Var(
		flagx.Func(func(s string) error {
			all = append(all, s)
			return nil
		}),
		"push", "push `value`",
	)

	flag.Parse()

	fmt.Println(all)

	// Output:
	// [a b]
}

// Shows an hex encoded parameter
func ExampleFunc_hex() {
	os.Args = []string{os.Args[0], "-hex=68656c6c6f"}
	flag.CommandLine = flag.NewFlagSet("test", flag.PanicOnError)

	// Destination of the decoded parameter value
	var bin []byte

	// A flag that apppends the given value to slice all
	flag.Var(
		flagx.Func(func(s string) (err error) {
			bin, err = hex.DecodeString(s)
			return
		}),
		"hex", "hex encoded `value`",
	)

	flag.Parse()

	fmt.Printf("%q", bin)

	// Output:
	// "hello"
}

func ExampleBoolFunc() {
	os.Args = []string{os.Args[0], "-rand", "-value=5"}
	flag.CommandLine = flag.NewFlagSet("test", flag.PanicOnError)

	var n int

	// A flag that set the value of n
	flag.IntVar(&n, "value", 0, "set given value")
	// A flag that set n to a random value
	flag.Var(
		flagx.BoolFunc(
			func(b bool) error {
				n = rand.Int()
				return nil
			},
		),
		"rand", "set random value",
	)

	flag.Parse()

	fmt.Println(n)

	// Output:
	// 5
}
