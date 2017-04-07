package flagx_test

import (
	"encoding/base64"
	"flag"
	"fmt"

	"github.com/dolmen-go/flagx"
)

func ExampleEncoded() {
	flags := flag.FlagSet{}

	var bin []byte
	// Bind parameter "-base64" to value bin above, with Base64 decoding
	flags.Var(flagx.Encoded(&bin, base64.RawStdEncoding), "base64", "hex string")

	if err := flags.Parse([]string{"-base64", "aGVsbG8K"}); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("%q\n", bin)

	// Output:
	// "hello\n"
}
