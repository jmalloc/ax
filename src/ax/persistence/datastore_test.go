package persistence_test

import (
	"context"

	. "github.com/jmalloc/ax/src/ax/persistence"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("WithDataStore / GetDataStore", func() {
	It("transports a datastore via the context", func() {
		in := &dataStore{}
		ctx := WithDataStore(context.Background(), in)

		out, ok := GetDataStore(ctx)

		Expect(ok).To(BeTrue())
		Expect(out).To(Equal(in))
	})
})

type dataStore struct {
	Begin func() Tx
}

func (ds *dataStore) BeginTx(context.Context) (Tx, error) {
	return ds.Begin(), nil
}
