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

// GetProto retrieves a protobuf-encoded message indexed with key from sequence of
// bkts, where each proceeding bucket is a subbucket of the previous one.
//
// It returns false if any of bkts do no exist or if the key in the last child
// bucket does not exist. It returns an error if one occurs in the process of
// unmarshaling binary data into m.
func GetProto(
	tx *bolt.Tx,
	key string,
	m proto.Message,
	bkts ...string,
) (bool, error) {
	pb := Get(tx, key, bkts...)
	if pb == nil {
		return false, nil
	}
	err := proto.Unmarshal(pb, m)
	if err != nil {
		return false, err
	}
	return true, nil
}

// PutProto places a protobuf-encoded message indexed with key to sequence of
// bkts, where each proceeding bucket is a subbucket of the previous one. If any
// of bkts don't exists, this function attempts to create them using MakeBtk
// function under the hood.
//
// It returns an error if bucket creation fails. It returns an error if one
// occurs in the process of marshalling m into binary.
func PutProto(
	tx *bolt.Tx,
	key string,
	m proto.Message,
	bkts ...string,
) error {
	pb, err := proto.Marshal(m)
	if err != nil {
		return err
	}

	return Put(tx, key, pb, bkts...)
}
