package format

import "fmt"

// Amount formats a cent amount as dollars.
func Amount(v int32) string {
	return fmt.Sprintf(
		"$%d.%02d",
		v/100,
		v%100,
	)
}
