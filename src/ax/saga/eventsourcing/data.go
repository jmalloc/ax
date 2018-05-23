package eventsourcing

import (
	"github.com/jmalloc/ax/src/ax"
	"github.com/jmalloc/ax/src/ax/saga"
)

// Data is a specialization of Data for sagas that use eventsourcing.
type Data interface {
	saga.Data

	// ApplyEvent updates the data to reflect the fact that an event has
	// occurred.
	//
	// It may panic if env.Message does not implement ax.Event.
	ApplyEvent(env ax.Envelope)
}
