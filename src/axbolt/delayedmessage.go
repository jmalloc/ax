package axbolt

import (
	"github.com/jmalloc/ax/src/ax/delayedmessage"
	boltdelayedmessage "github.com/jmalloc/ax/src/axbolt/delayedmessage"
)

// DelayedMessageRepository is an delayed message repository backed by a Bolt database.
var DelayedMessageRepository delayedmessage.Repository = boltdelayedmessage.Repository{}
