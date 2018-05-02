package persistence_test

import (
	"context"

	"github.com/jmalloc/ax/src/ax/internal/persistencetest"
	. "github.com/jmalloc/ax/src/ax/persistence"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("WithDataStore / GetDataStore", func() {
	It("transports a datastore via the context", func() {
		expected := &persistencetest.DataStoreMock{}
		ctx := WithDataStore(context.Background(), expected)

		ds, ok := GetDataStore(ctx)

		Expect(ok).To(BeTrue())
		Expect(ds).To(BeIdenticalTo(expected))
	})
})

type dataStore struct {
	Begin func() Tx
}

func (ds *dataStore) BeginTx(context.Context) (Tx, error) {
	return ds.Begin(), nil
}
