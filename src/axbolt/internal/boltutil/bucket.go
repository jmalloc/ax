package boltutil

import "github.com/boltdb/bolt"

// Has returns true if b contains k.
func Has(b *bolt.Bucket, k []byte) bool {
	return b.Get(k) != nil
}
