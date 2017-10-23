package flagnet_test

import (
	"flag"

	"github.com/dolmen-go/flagx/flagnet"
)

var x string
var _ flag.Getter = flagnet.HostPort(&x, -1)
