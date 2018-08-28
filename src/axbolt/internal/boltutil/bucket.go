package boltutil

import (
	"strings"

	"github.com/coreos/bbolt"
)

// Has returns true if b contains k.
func Has(b *bolt.Bucket, k []byte) bool {
	return b.Get(k) != nil
}

const (
	// bktPathSeparator is a delimeter that separates buckets
	// in the Bolt DB bucket path.
	bktPathSeparator = "/"
)

// GetBkt retrieves a Bolt DB bucket with a name of the last string of bkts.
// All preceding strings in bkts represent names the parent buckets of the last
// string (bucket name). If any bucket in bkts is nonexistant, nil is returned.
func GetBkt(tx *bolt.Tx, bkts ...string) *bolt.Bucket {
	var (
		b *bolt.Bucket
	)
	for _, s := range bkts {
		if b == nil {
			b = tx.Bucket([]byte(s))
			if b == nil {
				return nil
			}
		} else {
			b = b.Bucket([]byte(s))
			if b == nil {
				return nil
			}
		}
	}
	return b
}

// GetBktWithPath retrieves a Bolt DB bucket with the path of buckets separated
// with bktPathSeparator. This method retrieves the last bucket in the path. If
// any bucket in the path is nonexistant, nil is returned.
func GetBktWithPath(tx *bolt.Tx, path string) *bolt.Bucket {
	bkts := strings.Split(path, bktPathSeparator)
	return GetBkt(tx, bkts...)
}

// Get retrieves a value from the bucket which is the last in bkts and it has a
// key key. If any bucket in bkts is nonexistant nil is returned. If key does
// not exist in the last bucket of bkts, nil is returned.
func Get(tx *bolt.Tx, key string, bkts ...string) []byte {
	bkt := GetBkt(tx, bkts...)
	if bkt == nil {
		return nil
	}
	return bkt.Get([]byte(key))
}

// GetWithPath retrieves a value with key from the last bucket in the path. If
// any bucket in the path is nonexistant nil is returned. If key does not exist
// in the last bucket of the path, nil is returned.
func GetWithPath(tx *bolt.Tx, key, path string) []byte {
	bkt := GetBktWithPath(tx, path)
	if bkt == nil {
		return nil
	}
	return bkt.Get([]byte(key))
}

// MakeBkt creates Bolt DB buckets specified in bkts. Each subsequent bucket in
// bkts is created as a subbucket of preceding bucket. The last created bucket
// of bkts is returned under normal conditions. If any bucket specified in bkts
// already exists no action is taken on the bucket. If any Bolt DB error
// occurres in the process of creating the buckets, nil and the error are
// returned.
//
// Please note that this method trims any leading and trailing whitespaces
// (using strings.TrimSpace function) from all strings passed as bkts. It
// ignores creation of a bucket for any strings in bkts that result in zero
// length after trimming.
func MakeBkt(tx *bolt.Tx, bkts ...string) (*bolt.Bucket, error) {
	var (
		b   *bolt.Bucket
		err error
	)
	for _, s := range bkts {
		s = strings.TrimSpace(s)
		if len(s) > 0 {
			if b == nil {
				if b, err = tx.CreateBucketIfNotExists([]byte(s)); err != nil {
					return nil, err
				}
			} else {
				if b, err = b.CreateBucketIfNotExists([]byte(s)); err != nil {
					return nil, err
				}
			}
		}

	}
	return b, nil
}

// MakeBktWithPath creates Bolt DB buckets from path separated with bktPathSeparator
// For all usage cases please refer to function MakeBkt.
func MakeBktWithPath(tx *bolt.Tx, path string) (*bolt.Bucket, error) {
	bkts := strings.Split(path, bktPathSeparator)
	return MakeBkt(tx, bkts...)
}

// Put places a value val with key key in the bucket specified as a last string
// of bkts. Any nonexistant buckets in bkts will be created before the key and
// the value are placed. For any cases related to bucket creation of bkts,
// please refer to function MakeBkt.
func Put(tx *bolt.Tx, key string, val []byte, bkts ...string) error {
	bkt, err := MakeBkt(tx, bkts...)
	if err != nil {
		return err
	}
	return bkt.Put([]byte(key), []byte(val))
}

// PutWithPath places a value val with key key in the bucket specified as a last
// item in the path separated by bktPathSeparator. Any nonexistant buckets in the
// path will be created before the key and the value are placed. For any cases
// related to the creation of buckets in the path, please refer to function
// MakeBkt.
func PutWithPath(tx *bolt.Tx, key string, val []byte, path string) error {
	bkt, err := MakeBktWithPath(tx, path)
	if err != nil {
		return err
	}
	return bkt.Put([]byte(key), []byte(val))
}

// Delete deletes an entry with key key in the bucket specified as a last item
// in the bkts. If there are any nonexistant buckets in bkts the function
// returns nil.
func Delete(tx *bolt.Tx, key string, bkts ...string) error {
	if bkt := GetBkt(tx, bkts...); bkt != nil {
		return bkt.Delete([]byte(key))
	}
	return nil
}

// DeleteWithPath deletes an entry with key key in the bucket specified as a
// last in the separated by bktPathSeparator. If there are any nonexistant
// buckets in the path, the function returns nil.
func DeleteWithPath(tx *bolt.Tx, key string, path string) error {
	if bkt := GetBktWithPath(tx, path); bkt != nil {
		return bkt.Delete([]byte(key))
	}
	return nil
}
