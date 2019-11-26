package flagx_test

import (
	"encoding/hex"
	"flag"
	"fmt"

	"github.com/dolmen-go/flagx"
)

type hexString []byte

func (h hexString) MarshalText() (text []byte, err error) {
	hx := make([]byte, hex.EncodedLen(len(h)))
	_ = hex.Encode(hx, h)
	return hx, nil
}

func (h *hexString) UnmarshalText(text []byte) error {
	tmp := make(hexString, hex.DecodedLen(len(text)))
	_, err := hex.Decode(tmp, text)
	if err != nil {
		return err
	}
	*h = tmp
	return nil
}

func Example() {
	flags := flag.NewFlagSet("test", flag.PanicOnError) // Usually flag.CommandLine

	var hx hexString
	flags.Var(flagx.Text(&hx), "hex", "hex string")

	flags.Parse([]string{"-hex", "1234"})

	fmt.Printf("%d\n", hx)

	// Output:
	// [18 52]
}
