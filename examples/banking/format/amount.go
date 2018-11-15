package format

import "fmt"

// Amount formats a cent amount as dollars.
func Amount(v int32) string {
	f := "-$%d.%02d"
	if v < 0 {
		v = -v
		f = "-" + f
	}

	return fmt.Sprintf(
		"$%d.%02d",
		v/100,
		v%100,
	)
}
