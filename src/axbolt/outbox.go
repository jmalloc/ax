package axbolt

import (
	"github.com/jmalloc/ax/src/ax/outbox"
	boltoutbox "github.com/jmalloc/ax/src/axbolt/outbox"
)

// OutboxRepository is an outbox repository backed by a Bolt database.
var OutboxRepository outbox.Repository = boltoutbox.Repository{}
