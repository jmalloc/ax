package messagestoretests

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	// mysql driver blank import
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmalloc/ax/src/ax"
	"github.com/jmalloc/ax/src/ax/messagestore"
	"github.com/jmalloc/ax/src/ax/persistence"
	"github.com/jmalloc/ax/src/axtest"
	"github.com/jmalloc/ax/src/axtest/testmessages"
	g "github.com/onsi/ginkgo"
	m "github.com/onsi/gomega"
)

// DumpSQLTable dumps MySQL table into the standard output
func DumpSQLTable(t string) error {
	dsn := os.Getenv("AX_MYSQL_DSN")
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return err
	}

	rows, err := db.Query(fmt.Sprintf("SELECT * FROM %s;", t))
	if err != nil {
		return err
	}
	columns, _ := rows.Columns()
	count := len(columns)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)

	for rows.Next() {
		fmt.Println("----------")

		for i := range columns {
			valuePtrs[i] = &values[i]
		}

		rows.Scan(valuePtrs...)
		for i, col := range columns {
			var v interface{}
			val := values[i]
			b, ok := val.([]byte)

			if ok {
				v = string(b)
			} else {
				v = val
			}
			fmt.Println(col, ": ", v)
		}
	}
	return nil
}

// MessageStoreSuite returns a test suite for implementations of messagestore.GloballyOrderedStore,
func MessageStoreSuite(
	getStore func() persistence.DataStore,
	getMessageStore func() messagestore.GloballyOrderedStore,
) func() {
	return func() {
		var (
			causationID, correlationID ax.MessageID
			store                      persistence.DataStore
			msgStore                   messagestore.GloballyOrderedStore
			ctx                        context.Context
			cancel                     func()
			m1, m2                     ax.Envelope
			t1, t2, t3, t4             time.Time
			s1, s2                     string
		)

		g.BeforeEach(func() {
			store = getStore()
			msgStore = getMessageStore()
			var fn func()
			ctx, fn = context.WithTimeout(context.Background(), 15*time.Second)
			cancel = fn // defeat go vet warning about unused cancel func

			causationID = ax.GenerateMessageID()
			correlationID = ax.GenerateMessageID()

			t1 = time.Now()
			t2 = time.Now()
			t3 = time.Now()
			t4 = time.Now()
			s1 = "<stream1>"
			s2 = "<stream2>"

			m1 = ax.Envelope{
				MessageID:     ax.GenerateMessageID(),
				CausationID:   causationID,
				CorrelationID: correlationID,
				CreatedAt:     t1,
				SendAt:        t2,
				Message: &testmessages.Command{
					Value: "<foo>",
				},
			}
			m2 = ax.Envelope{
				MessageID:     ax.GenerateMessageID(),
				CausationID:   causationID,
				CorrelationID: correlationID,
				CreatedAt:     t3,
				SendAt:        t4,
				Message: &testmessages.Event{
					Value: "<bar>",
				},
			}
		})

		g.AfterEach(func() {
			cancel()
		})

		g.Describe("AppendMessages", func() {
			g.Context("when the offset is the next unused offset in the stream", func() {
				g.It("returns no error", func() {
					tx, com, err := store.BeginTx(ctx)
					m.Expect(err).ShouldNot(m.HaveOccurred())
					defer com.Rollback()

					offset := uint64(0)

					err = msgStore.AppendMessages(
						ctx,
						tx,
						s1,
						offset,
						[]ax.Envelope{m1, m2},
					)
					m.Expect(err).ShouldNot(m.HaveOccurred())

					err = msgStore.AppendMessages(
						ctx,
						tx,
						s2,
						offset,
						[]ax.Envelope{m1, m2},
					)
					m.Expect(err).ShouldNot(m.HaveOccurred())

					err = com.Commit()
					m.Expect(err).ShouldNot(m.HaveOccurred())
				})
			})
			g.Context("when the offset is not the next unused offset in the stream", func() {
				g.It("returns an error", func() {
					tx, com, err := store.BeginTx(ctx)
					m.Expect(err).ShouldNot(m.HaveOccurred())
					defer com.Rollback()

					offset := uint64(999)
					err = msgStore.AppendMessages(
						ctx,
						tx,
						s1,
						offset,
						[]ax.Envelope{m1, m2},
					)
					m.Expect(err).Should(m.HaveOccurred())

					err = com.Commit()
					m.Expect(err).ShouldNot(m.HaveOccurred())
				})
			})
		})

		g.Describe("OpenStream", func() {
			g.BeforeEach(func() {
				tx, com, err := store.BeginTx(ctx)
				m.Expect(err).ShouldNot(m.HaveOccurred())
				defer com.Rollback()

				offset := uint64(0)

				err = msgStore.AppendMessages(
					ctx,
					tx,
					s1,
					offset,
					[]ax.Envelope{m1, m2},
				)
				m.Expect(err).ShouldNot(m.HaveOccurred())

				err = msgStore.AppendMessages(
					ctx,
					tx,
					s2,
					offset,
					[]ax.Envelope{m1, m2},
				)
				m.Expect(err).ShouldNot(m.HaveOccurred())

				err = com.Commit()
				m.Expect(err).ShouldNot(m.HaveOccurred())
			})
			g.Context("when stream exists", func() {
				g.It("returns true", func() {
					offset := uint64(0)
					_, ok, err := msgStore.OpenStream(
						ctx,
						store,
						s1,
						offset,
					)
					m.Expect(err).ShouldNot(m.HaveOccurred())
					m.Expect(ok).Should(m.BeTrue())
				})
				g.It("returns the stream", func() {
					offset := uint64(0)
					s, _, err := msgStore.OpenStream(
						ctx,
						store,
						s1,
						offset,
					)
					m.Expect(err).ShouldNot(m.HaveOccurred())
					m.Expect(s).ShouldNot(m.BeNil())
				})
			})
			g.Context("Stream.Next", func() {
				g.Context("when next message is available", func() {
					g.It("advances the stream and returns nil", func() {
						offset := uint64(0)
						s, ok, err := msgStore.OpenStream(
							ctx,
							store,
							s1,
							offset,
						)
						m.Expect(err).ShouldNot(m.HaveOccurred())
						m.Expect(ok).Should(m.BeTrue())

						err = s.Next(ctx)
						m.Expect(err).ShouldNot(m.HaveOccurred())
					})
				})
				g.Context("when context is canceled", func() {
					g.It("returns context.Canceled error", func() {
						offset := uint64(0)
						s, ok, err := msgStore.OpenStream(
							ctx,
							store,
							s1,
							offset,
						)
						m.Expect(err).ShouldNot(m.HaveOccurred())
						m.Expect(ok).Should(m.BeTrue())

						cancel()

						err = s.Next(ctx)
						m.Expect(err).Should(m.MatchError(context.Canceled))
					})
				})
				g.Context("when next message is unavailable", func() {
					g.It("blocks indefinitely", func() {
						offset := uint64(0)
						errNotify := make(chan error)
						s, ok, err := msgStore.OpenStream(
							ctx,
							store,
							s1,
							offset,
						)
						m.Expect(err).ShouldNot(m.HaveOccurred())
						m.Expect(ok).Should(m.BeTrue())

						go func() {
							for {
								errNotify <- s.Next(ctx)
							}
						}()
						// m1
						m.Eventually(errNotify).Should(m.Receive(m.Succeed()))
						// m2
						m.Eventually(errNotify).Should(m.Receive(m.Succeed()))
						// no other messages
						m.Consistently(errNotify).ShouldNot(m.Receive())
					})
				})
			})
			g.Context("Stream.Offset", func() {
				g.It("returns the current offset of the stream", func() {
					offset := uint64(0)
					s, ok, err := msgStore.OpenStream(
						ctx,
						store,
						s1,
						offset,
					)
					m.Expect(err).ShouldNot(m.HaveOccurred())
					m.Expect(ok).Should(m.BeTrue())

					err = s.Next(ctx)
					m.Expect(err).ShouldNot(m.HaveOccurred())

					o, err := s.Offset()
					m.Expect(err).ShouldNot(m.HaveOccurred())
					m.Expect(o).Should(m.BeNumerically("==", 0))

					err = s.Next(ctx)
					m.Expect(err).ShouldNot(m.HaveOccurred())

					o, err = s.Offset()
					m.Expect(err).ShouldNot(m.HaveOccurred())
					m.Expect(o).Should(m.BeNumerically("==", 1))
				})
			})
			g.Context("Stream.Get", func() {
				g.It("returns the message at the current offset of the stream", func() {
					offset := uint64(0)
					s, ok, err := msgStore.OpenStream(
						ctx,
						store,
						s1,
						offset,
					)
					m.Expect(err).ShouldNot(m.HaveOccurred())
					m.Expect(ok).Should(m.BeTrue())

					err = s.Next(ctx)
					m.Expect(err).ShouldNot(m.HaveOccurred())

					env, err := s.Get(ctx)
					m.Expect(err).ShouldNot(m.HaveOccurred())
					m.Expect(axtest.EnvelopesEqual(env, m1)).Should(m.BeTrue())

					err = s.Next(ctx)
					m.Expect(err).ShouldNot(m.HaveOccurred())

					env, err = s.Get(ctx)
					m.Expect(err).ShouldNot(m.HaveOccurred())
					m.Expect(axtest.EnvelopesEqual(env, m2)).Should(m.BeTrue())
				})
			})
		})

		g.Describe("OpenGlobal", func() {
			g.BeforeEach(func() {
				tx, com, err := store.BeginTx(ctx)
				m.Expect(err).ShouldNot(m.HaveOccurred())
				defer com.Rollback()

				offset := uint64(0)

				err = msgStore.AppendMessages(
					ctx,
					tx,
					s1,
					offset,
					[]ax.Envelope{m1, m2},
				)
				m.Expect(err).ShouldNot(m.HaveOccurred())

				err = com.Commit()
				m.Expect(err).ShouldNot(m.HaveOccurred())

				tx, com, err = store.BeginTx(ctx)
				m.Expect(err).ShouldNot(m.HaveOccurred())
				defer com.Rollback()

				err = msgStore.AppendMessages(
					ctx,
					tx,
					s2,
					offset,
					[]ax.Envelope{m1, m2},
				)
				m.Expect(err).ShouldNot(m.HaveOccurred())

				err = com.Commit()
				m.Expect(err).ShouldNot(m.HaveOccurred())
			})
			g.Context("when global stream exists", func() {
				g.It("returns no error", func() {
					offset := uint64(0)
					_, err := msgStore.OpenGlobal(
						ctx,
						store,
						offset,
					)
					m.Expect(err).ShouldNot(m.HaveOccurred())
				})
			})
			g.Context("Stream.Next", func() {
				g.Context("when next message is available", func() {
					g.It("advances the stream and returns nil", func() {
						offset := uint64(0)
						s, err := msgStore.OpenGlobal(
							ctx,
							store,
							offset,
						)
						m.Expect(err).ShouldNot(m.HaveOccurred())

						err = s.Next(ctx)
						m.Expect(err).ShouldNot(m.HaveOccurred())
					})
				})
				g.Context("when context is canceled", func() {
					g.It("returns context.Canceled error", func() {
						offset := uint64(0)
						s, err := msgStore.OpenGlobal(
							ctx,
							store,
							offset,
						)
						m.Expect(err).ShouldNot(m.HaveOccurred())

						cancel()

						err = s.Next(ctx)
						m.Expect(err).Should(m.MatchError(context.Canceled))
					})
				})
				g.Context("when next message is unavailable", func() {
					g.It("blocks indefinitely", func() {
						offset := uint64(0)
						errNotify := make(chan error)
						s, err := msgStore.OpenGlobal(
							ctx,
							store,
							offset,
						)
						m.Expect(err).ShouldNot(m.HaveOccurred())

						go func() {
							for {
								errNotify <- s.Next(ctx)
							}
						}()
						// m1
						m.Eventually(errNotify).Should(m.Receive(m.Succeed()))
						// m2
						m.Eventually(errNotify).Should(m.Receive(m.Succeed()))
						// m3
						m.Eventually(errNotify).Should(m.Receive(m.Succeed()))
						// m4
						m.Eventually(errNotify).Should(m.Receive(m.Succeed()))
						// no other messages
						m.Consistently(errNotify).ShouldNot(m.Receive())
					})
				})
			})
			g.Context("Stream.Offset", func() {
				g.It("returns the current offset of the stream", func() {
					offset := uint64(0)
					s, err := msgStore.OpenGlobal(
						ctx,
						store,
						offset,
					)
					m.Expect(err).ShouldNot(m.HaveOccurred())

					err = s.Next(ctx)
					m.Expect(err).ShouldNot(m.HaveOccurred())

					o, err := s.Offset()
					m.Expect(err).ShouldNot(m.HaveOccurred())
					m.Expect(o).Should(m.BeNumerically("==", 0))

					err = s.Next(ctx)
					m.Expect(err).ShouldNot(m.HaveOccurred())

					o, err = s.Offset()
					m.Expect(err).ShouldNot(m.HaveOccurred())
					m.Expect(o).Should(m.BeNumerically("==", 1))
				})
			})
			g.Context("Stream.Get", func() {
				g.It("returns the message at the current offset of the stream", func() {
					offset := uint64(0)
					s, err := msgStore.OpenGlobal(
						ctx,
						store,
						offset,
					)
					m.Expect(err).ShouldNot(m.HaveOccurred())

					err = s.Next(ctx)
					m.Expect(err).ShouldNot(m.HaveOccurred())

					env, err := s.Get(ctx)
					m.Expect(err).ShouldNot(m.HaveOccurred())
					m.Expect(axtest.EnvelopesEqual(env, m1)).Should(m.BeTrue())

					err = s.Next(ctx)
					m.Expect(err).ShouldNot(m.HaveOccurred())

					env, err = s.Get(ctx)
					m.Expect(err).ShouldNot(m.HaveOccurred())
					err = DumpSQLTable("ax_messagestore_message")
					m.Expect(err).ShouldNot(m.HaveOccurred())
					m.Expect(axtest.EnvelopesEqual(env, m2)).Should(m.BeTrue())
				})
			})
		})
	}
}
