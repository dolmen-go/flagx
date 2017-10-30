package flagx

import (
	"time"
)

// Time allows to define a flag that feeds a time.Time value
type Time struct {
	Time     *time.Time
	Format   string // See time.ParseTime
	Location *time.Location
}

func (t Time) String() string {
	if t.Time == nil || t.Time.IsZero() {
		// When called by flag.isZeroValue
		return ""
	}
	return t.Time.Format(t.Format)
}

func (t Time) Set(s string) (err error) {
	if t.Location != nil {
		*t.Time, err = time.ParseInLocation(t.Format, s, t.Location)
	} else {
		*t.Time, err = time.Parse(t.Format, s)
	}
	return
}

func (t Time) Get() interface{} {
	return *t.Time
}
