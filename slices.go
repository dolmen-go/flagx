package flagx

import (
	"fmt"
	"strconv"
	"strings"
)

type IntSlice struct {
	Slice *[]int
}

func (is IntSlice) String() string {
	if is.Slice == nil {
		// When called by flag.isZeroValue
		return ""
	}
	if *is.Slice == nil {
		return ""
	}
	str := make([]string, len(*is.Slice))
	for i := range str {
		str[i] = strconv.Itoa((*is.Slice)[i])
	}
	return strings.Join(str, ",")
}

func (is IntSlice) Set(s string) (err error) {
	str := strings.Split(s, ",")
	for _, s := range str {
		s = strings.TrimSpace(s)
		n, err := strconv.ParseInt(s, 0, 0)
		if err != nil {
			return fmt.Errorf("%q: %s", s, err)
		}
		*is.Slice = append(*is.Slice, int(n))
	}
	return nil
}

func (is IntSlice) Get() interface{} {
	return *is.Slice
}
