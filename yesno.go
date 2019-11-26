package flagx

import "strconv"

// YesNo returns a flag.Value for a boolean value, but accepting "yes"/"y"/"no"/"n" in addition to strconv.ParseBool values.
func YesNo(b *bool) Value {
	_ = *b // Check not nil
	return yesNo{b}
}

type yesNo struct {
	b *bool
}

func (yn yesNo) Set(s string) error {
	switch s {
	case "yes", "y", "true", "1":
		*yn.b = true
	case "no", "n", "false", "0":
		*yn.b = false
	default:
		v, err := strconv.ParseBool(s)

		if err == nil {
			*yn.b = v
		}
		return err

	}
	return nil
}

func (yn yesNo) Get() interface{} { return *yn.b }

func (yn yesNo) String() string { return strconv.FormatBool(*yn.b) }

func (yn yesNo) IsBoolFlag() bool { return true }
