package messagestore

import (
	"context"
	"database/sql"
)

// fetchColumns is the ordered list of columns that must be SELECTed by fetchers.
const fetchColumns = `message_id,
					  causation_id,
					  correlation_id,
					  time,
					  content_type,
					  data`

// Fetcher is an interface for fetching rows from the message store
type Fetcher interface {
	// FetchRows fetches the n rows beginning at the given offset.
	FetchRows(ctx context.Context, offset, n uint64) (*sql.Rows, error)
}

// StreamFetcher is a fetcher that fetches rows for a specific stream.
type StreamFetcher struct {
	DB       *sql.DB
	StreamID int64
}

// FetchRows fetches the n rows beginning at the given offset.
func (f *StreamFetcher) FetchRows(ctx context.Context, offset, n uint64) (*sql.Rows, error) {
	return f.DB.QueryContext(
		ctx,
		`SELECT `+fetchColumns+`
		FROM ax_messagestore_message
		WHERE stream_id = ?
		AND stream_offset >= ?
		ORDER BY stream_offset
		LIMIT ?`,
		f.StreamID,
		offset,
		n,
	)
}

// GlobalFetcher is a fetcher that fetches rows for the entire store
type GlobalFetcher struct {
	DB *sql.DB
}

// FetchRows fetches the n rows beginning at the given offset.
func (f *GlobalFetcher) FetchRows(ctx context.Context, offset, n uint64) (*sql.Rows, error) {
	return f.DB.QueryContext(
		ctx,
		`SELECT `+fetchColumns+`
		FROM ax_messagestore_message
		WHERE global_offset >= ?
		ORDER BY global_offset
		LIMIT ?`,
		offset,
		n,
	)
}
