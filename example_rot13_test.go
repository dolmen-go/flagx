package flagx_test

import (
	"flag"
	"fmt"

	"github.com/dolmen-go/flagx"
)

var rot13replacer [256]byte

func init() {
	for i := 0; i < 256; i++ {
		rot13replacer[i] = byte(i)
	}
	const lower = 'a' - 'A'
	for i := 'A'; i <= 'L'; i++ {
		rot13replacer[i] = byte(i + 13)
		rot13replacer[i+13] = byte(i)
		rot13replacer[i+lower] = byte(i + lower + 13)
		rot13replacer[i+lower+13] = byte(i + lower)
	}
}

type rot13 struct{}

func (rot13) EncodedLen(n int) int { return n }

func (rot13) Encode(dst, src []byte) {
	if len(src) == 0 {
		return
	}
	_ = dst[len(src)-1]
	for i := 0; i < len(src); i++ {
		dst[i] = rot13replacer[src[i]]
	}
}

// EncodeToString returns the rot13 encoding of src.
func (rot13) EncodeToString(src []byte) string {
	dst := make([]byte, len(src))
	rot13{}.Encode(dst, src)
	return string(dst)
}

func (rot13) DecodedLen(n int) int { return n }

func (rot13) Decode(dst, src []byte) (int, error) {
	rot13{}.Encode(dst, src)
	return len(src), nil
}

func (rot13) DecodeString(src string) ([]byte, error) {
	dst := make([]byte, len(src))
	for i := 0; i < len(src); i++ {
		dst[i] = rot13replacer[src[i]]
	}
	return dst, nil
}

func ExampleEncoded_rot13() {
	flags := flag.NewFlagSet("test", flag.PanicOnError) // Usually flag.CommandLine

	var bin []byte
	// Bind parameter "-password" to value bin above, with ROT13 decoding
	flags.Var(flagx.Encoded(&bin, rot13{}), "password", "")

	flags.Parse([]string{"-password", "frperg"})

	fmt.Printf("Decoded: %q\n", bin)

	// Output:
	// Decoded: "secret"
}
