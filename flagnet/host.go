package flagnet

import (
	"fmt"
	"net"
)

func Host(value *string) Value {
	_ = *value // Trigger a panic if value == nil
	return &host{value}
}

type host struct {
	Host *string
}

func (h host) String() string {
	return *h.Host
}

func (h host) Set(str string) (err error) {
	_, err = net.LookupHost(str)
	if err != nil {
		return fmt.Errorf("%q: %s", str, err)
	}
	*h.Host = str
	return nil
}

func (h host) Get() interface{} {
	return *h.Host
}
