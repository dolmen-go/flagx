package flagx_test

import (
	"encoding/hex"
	"flag"
	"fmt"
	"math/rand"

	"github.com/dolmen-go/flagx"
)

func ExampleFunc() {
	flags := flag.NewFlagSet("test", flag.PanicOnError) // Usually flag.CommandLine

	var all []string

	push := func(s string) error {
		all = append(all, s)
		return nil
	}

	// A flag that apppends the given value to slice all
	flags.Var(
		flagx.Func(push),
		"push", "push `value`",
	)

	flags.Parse([]string{"-push=a", "-push=b"})

	fmt.Println(all)

	// Output:
	// [a b]
}

// Shows an hex encoded parameter
func ExampleFunc_hex() {
	flags := flag.NewFlagSet("test", flag.PanicOnError) // Usually flag.CommandLine

	// Destination of the decoded parameter value
	var bin []byte

	// A flag that decodes value from hexadecimal
	flags.Var(
		flagx.Func(func(s string) (err error) {
			bin, err = hex.DecodeString(s)
			return
		}),
		"hex", "hex encoded `value`",
	)

	flags.Parse([]string{"-hex=68656c6c6f"})

	fmt.Printf("%q", bin)

	// Output:
	// "hello"
}

func ExampleBoolFunc() {
	flags := flag.NewFlagSet("test", flag.PanicOnError) // Usually flag.CommandLine

	var n int

	// A flag that set the value of n
	flags.IntVar(&n, "value", 0, "set given value")
	// A flag that set n to a random value
	flags.Var(
		flagx.BoolFunc(
			func(b bool) error {
				n = rand.Int()
				return nil
			},
		),
		"rand", "set random value",
	)

	flags.Parse([]string{"-rand", "-value=5"})

	fmt.Println(n)

	// Output:
	// 5
}
