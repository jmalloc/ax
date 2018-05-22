package saga

// KeySet is a set of "mapping keys" that are associated with a saga instance.
type KeySet map[string]struct{}

// NewKeySet returns a key set containing the keys in k.
func NewKeySet(k ...string) KeySet {
	s := make(KeySet, len(k))

	for _, key := range k {
		s[key] = struct{}{}
	}

	return s
}
