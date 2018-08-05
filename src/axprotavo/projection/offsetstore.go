package projection

import (
	"context"
	"fmt"

	"github.com/jmalloc/ax/src/ax/persistence"
	protavopersistence "github.com/jmalloc/ax/src/axprotavo/persistence"
	"github.com/jmalloc/protavo/src/protavo"
	"github.com/jmalloc/protavo/src/protavo/document"
)

// OffsetStore is a Protavo-backed implementation of Ax's projection.OffsetStore
// interface.
type OffsetStore struct{}

// LoadOffset returns the offset at which a consumer should resume
// reading from the stream.
//
// pk is the projector's persistence key.
func (OffsetStore) LoadOffset(
	ctx context.Context,
	ds persistence.DataStore,
	pk string,
) (uint64, error) {
	db := protavopersistence.ExtractDB(ds)

	id := persistenceKeyToDocID(pk)
	doc, ok, err := db.Load(ctx, id)
	if !ok || err != nil {
		return 0, err
	}

	return doc.Revision, nil
}

// IncrementOffset increments the offset at which a consumer should resume
// reading from the stream by one.
//
// pk is the projector's persitence key. c is the offset that is currently
// stored, as returned by LoadOffset(). If c is not the offset that is
// currently stored, the increment fails and a non-nil error is returned.
func (s OffsetStore) IncrementOffset(
	ctx context.Context,
	ptx persistence.Tx,
	pk string,
	c uint64,
) error {
	tx := protavopersistence.ExtractTx(ptx)

	// note that we just use the document revision to store the offset. we don't
	// need any additional content.
	doc := &document.Document{
		ID:       persistenceKeyToDocID(pk),
		Revision: c,
		Content:  document.StringContent(""), // TODO: remove when https://github.com/jmalloc/protavo/issues/6 is complete
	}

	op := protavo.Save(doc)
	op.ExecuteInWriteTx(ctx, tx)

	err := op.Err()

	if protavo.IsOptimisticLockError(err) {
		// TODO: use OCC error https://github.com/jmalloc/ax/issues/93
		return fmt.Errorf(
			"can not increment projection offset for persistence key %s, offset %d is not the current offset",
			pk,
			c,
		)
	}

	return err
}

func persistenceKeyToDocID(pk string) string {
	return "persistence-offset-" + pk
}
