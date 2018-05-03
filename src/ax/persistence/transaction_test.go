package persistence_test

import (
	"context"

	. "github.com/jmalloc/ax/src/ax/persistence"
	"github.com/jmalloc/ax/src/internal/persistencetest"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("WithTransaction / GetTransaction", func() {
	It("transports a transaction via the context", func() {
		expected := &persistencetest.TxMock{}
		ctx := WithTx(context.Background(), expected)

		tx, ok := GetTx(ctx)

		Expect(ok).To(BeTrue())
		Expect(tx).To(BeIdenticalTo(expected))
	})
})

var _ = Describe("GetOrBeginTx", func() {
	Context("when the context already contains a transaction", func() {
		tx := &persistencetest.TxMock{}
		ctx := WithTx(context.Background(), tx)

		It("returns the transaction", func() {
			t, _, err := GetOrBeginTx(ctx)

			Expect(err).ShouldNot(HaveOccurred())
			Expect(t).To(BeIdenticalTo(tx))
		})

		It("returns a no-op committer", func() {
			_, c, err := GetOrBeginTx(ctx)
			Expect(err).ShouldNot(HaveOccurred())

			err = c.Commit()
			Expect(err).ShouldNot(HaveOccurred())
		})
	})

	Context("when the context does not already contain a transaction", func() {
		tx := &persistencetest.TxMock{}
		com := &persistencetest.CommitterMock{}
		ds := &persistencetest.DataStoreMock{
			BeginTxFunc: func(context.Context) (Tx, Committer, error) {
				return tx, com, nil
			},
		}

		Context("when the context contains a data store", func() {
			ctx := WithDataStore(context.Background(), ds)

			It("starts a new transaction using the data store", func() {
				t, _, err := GetOrBeginTx(ctx)

				Expect(err).ShouldNot(HaveOccurred())
				Expect(t).To(BeIdenticalTo(tx))
			})

			It("returns the committer returned by the data store", func() {
				_, c, err := GetOrBeginTx(ctx)

				Expect(err).ShouldNot(HaveOccurred())
				Expect(c).To(BeIdenticalTo(com))
			})
		})

		Context("when the context does not contains a data store", func() {
			It("returns an error", func() {
				_, _, err := GetOrBeginTx(context.Background())

				Expect(err).Should(HaveOccurred())
			})
		})
	})
})
