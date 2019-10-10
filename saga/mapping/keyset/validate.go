package keyset

import (
	"errors"
	"fmt"
)

// Validate returns a copy of ks with duplicate keys removed.
// It returns an error of any of the keys is the empty string.
func Validate(ks []string) ([]string, error) {
	dedup := make([]string, 0, len(ks))
	seen := make(map[string]struct{}, len(ks))

	for _, k := range ks {
		if k == "" {
			return nil, errors.New("mapping keys must not be empty")
		}

		if _, ok := seen[k]; ok {
			return nil, fmt.Errorf("the mapping key %s is repeated in the key set", k)
		}

		seen[k] = struct{}{}
		dedup = append(dedup, k)
	}

	return dedup, nil
}
