package persistence_test

import (
	"context"

	. "github.com/jmalloc/ax/src/ax/persistence"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("WithTransaction / GetTransaction", func() {
	It("transports a transaction via the context", func() {
		in := &transaction{}
		ctx := WithTx(context.Background(), in)

		out, ok := GetTx(ctx)

		Expect(ok).To(BeTrue())
		Expect(out).To(Equal(in))
	})
})

var _ = Describe("GetOrBeginTx", func() {
	Context("when the context already contains a transaction", func() {
		tx := &transaction{
			Underlying: "<underlying tx>",
		}
		ctx := WithTx(context.Background(), tx)

		It("returns a transaction that exposes the same underlying transaction", func() {
			t, err := GetOrBeginTx(ctx)

			Expect(err).ShouldNot(HaveOccurred())
			Expect(t.UnderlyingTx()).To(Equal(tx.Underlying))
		})

		It("does not allow the transaction to be committed", func() {
			t, _ := GetOrBeginTx(ctx)
			t.Commit()

			Expect(tx.CommitCalled).To(BeFalse())
		})

		It("does not allow the transaction to be rolled back", func() {
			t, _ := GetOrBeginTx(ctx)
			t.Rollback()

			Expect(tx.RollbackCalled).To(BeFalse())
		})
	})

	Context("when the context does not already contain a transaction", func() {
		tx := &transaction{}
		ds := &dataStore{
			Begin: func() Tx {
				return tx
			},
		}

		Context("when the context contains a data store", func() {
			ctx := WithDataStore(context.Background(), ds)

			It("starts a new transaction using the data store", func() {
				t, err := GetOrBeginTx(ctx)

				Expect(err).ShouldNot(HaveOccurred())
				Expect(t).To(Equal(tx))
			})

			It("allows the transaction to be committed", func() {
				t, _ := GetOrBeginTx(ctx)
				t.Commit()

				Expect(tx.CommitCalled).To(BeTrue())
			})

			It("allows the transaction to be rolled back", func() {
				t, _ := GetOrBeginTx(ctx)
				t.Rollback()

				Expect(tx.RollbackCalled).To(BeTrue())
			})
		})
	})

	Context("when the context does not contains a data store", func() {
		It("returns an error ", func() {
			_, err := GetOrBeginTx(context.Background())

			Expect(err).Should(HaveOccurred())
		})
	})
})

var _ = Describe("BeginTx", func() {

})

type transaction struct {
	CommitCalled   bool
	RollbackCalled bool
	Underlying     interface{}
}

func (tx *transaction) Commit() error {
	tx.CommitCalled = true
	return nil
}

func (tx *transaction) Rollback() error {
	tx.RollbackCalled = true
	return nil
}

func (tx *transaction) UnderlyingTx() interface{} {
	return tx.Underlying
}
