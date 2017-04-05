package flagtext_test

import (
	"encoding/hex"
	"flag"
	"fmt"

	"github.com/dolmen-go/flagtext"
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
	flags := flag.FlagSet{}

	var hx hexString
	flags.Var(flagtext.Text(&hx), "hex", "hex string")

	if err := flags.Parse([]string{"-hex", "1234"}); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("%d\n", hx)

	// Output:
	// [18 52]
}
