package flagx_test

import (
	"encoding/hex"
	"flag"
	"fmt"

	"github.com/dolmen-go/flagx"
)

type hexEncoding struct{}

func (hexEncoding) DecodeString(s string) ([]byte, error) {
	return hex.DecodeString(s)
}

func (hexEncoding) EncodeToString(src []byte) string {
	return hex.EncodeToString(src)
}

func ExampleEncoded_hex() {
	flags := flag.NewFlagSet("test", flag.PanicOnError) // Usually flag.CommandLine

	var bin []byte
	// Bind parameter "-hex" to value bin above, with hex decoding
	flags.Var(flagx.Encoded(&bin, hexEncoding{}), "hex", "hex string")

	flags.Parse([]string{"-hex", "68656c6c6f"})

	fmt.Printf("%q\n", bin)

	// Output:
	// "hello"
}
