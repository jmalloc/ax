// Package messagestore contains interfaces used by Ax to read and write from
// persisted streams of messages.
//
// The message store is the fundamental persistence type used by eventsourced
// sagas, though the interface is not restricted to storing events.
package messagestore
