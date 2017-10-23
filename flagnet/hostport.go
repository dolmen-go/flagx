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
	if h.HostPort == nil {
		// When called by flag.isZeroValue
		return ""
	}
	return *h.HostPort
}

func (h hostPort) Set(str string) (err error) {
	host, _, err := net.SplitHostPort(str)
	if err != nil {
		if h.DefaultPort < 0 {
			return fmt.Errorf("%q: %s", str, err)
		}
		// If the port is missing, append it
		str = str + ":" + strconv.Itoa(h.DefaultPort)
		host, _, err = net.SplitHostPort(str)
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
