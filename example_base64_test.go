package flagx_test

import (
	"encoding/base64"
	"flag"
	"fmt"

	"github.com/dolmen-go/flagx"
)

func ExampleEncoded_base64() {
	flags := flag.NewFlagSet("test", flag.PanicOnError) // Usually flag.CommandLine

	var bin []byte
	// Bind parameter "-base64" to value bin above, with Base64 decoding
	flags.Var(flagx.Encoded(&bin, base64.RawStdEncoding), "base64", "hex string")

	flags.Parse([]string{"-base64", "aGVsbG8K"})

	fmt.Printf("%q\n", bin)

	// Output:
	// "hello\n"
}
