package flagnet

import (
	"fmt"
	"net"
	"strconv"
)

func HostPort(value *string, defaultPort int) Value {
	return &hostPort{value, defaultPort}
}

type hostPort struct {
	HostPort    *string
	DefaultPort int
}

func (h hostPort) String() string {
	return *h.HostPort
}

func (h hostPort) Set(str string) (err error) {
	host, _, err := net.SplitHostPort(str)
	if err != nil {
		if h.DefaultPort >= 0 {
			host, _, err = net.SplitHostPort(str + ":" + strconv.Itoa(h.DefaultPort))
		}
		if err != nil {
			return fmt.Errorf("%q: %s", str, err)
		}
	}

	if _, err = net.LookupHost(host); err != nil {
		return fmt.Errorf("%q: %s", host, err)
	}

	*h.HostPort = str

	return nil
}

func (h hostPort) Get() interface{} {
	return *h.HostPort
}
