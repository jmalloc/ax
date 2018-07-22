package boltutil

import (
	"github.com/coreos/bbolt"
	"github.com/golang/protobuf/proto"
)

// MarshalProto stores m in b with key k.
func MarshalProto(b *bolt.Bucket, k []byte, m proto.Message) error {
	buf, err := proto.Marshal(m)
	if err != nil {
		return err
	}

	return b.Put(k, buf)
}

// UnmarshalProto populates m from the value stored in b with key k.
// It returns false if k is empty.
func UnmarshalProto(b *bolt.Bucket, k []byte, m proto.Message) (bool, error) {
	buf := b.Get(k)
	if buf == nil {
		return false, nil
	}

	return true, proto.Unmarshal(buf, m)
}
