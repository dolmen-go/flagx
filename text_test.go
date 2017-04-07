package flagx

import (
	"testing"

	"bytes"
	"encoding/hex"
	"flag"
)

var _ flag.Getter = (*textValue)(nil)

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

func TestVarText(t *testing.T) {
	flags := flag.FlagSet{}
	var hx hexString
	flags.Var(Text(&hx), "hex", "hex string")
	if err := flags.Parse([]string{"-hex", "1234"}); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if !bytes.Equal(hx, []byte{0x12, 0x34}) {
		t.Fatalf("got: %#v", hx)
	}
}
