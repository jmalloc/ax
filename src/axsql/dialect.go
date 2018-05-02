package axsql

import (
	"context"
	"database/sql"

	"github.com/jmalloc/ax/src/ax"
	"github.com/jmalloc/ax/src/ax/bus"
)

// Dialect is an interface that hides the difference between various SQL
// dialects.
type Dialect interface {
	TxOptions() *sql.TxOptions

	SelectOutboxExists(context.Context, *sql.DB, ax.MessageID) (bool, error)
	SelectOutboxEnvelopes(context.Context, *sql.DB, ax.MessageID) ([]bus.OutboundEnvelope, error)
	InsertOutbox(context.Context, *sql.Tx, ax.MessageID) error
	InsertOutboxEnvelope(context.Context, *sql.Tx, bus.OutboundEnvelope) error
	DeleteOutboxEnvelope(context.Context, *sql.Tx, ax.MessageID) error
}
