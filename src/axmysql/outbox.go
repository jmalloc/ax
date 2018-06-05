package axmysql

import (
	"github.com/jmalloc/ax/src/ax/outbox"
	mysqloutbox "github.com/jmalloc/ax/src/axmysql/outbox"
)

// OutboxRepository is an outbox repository backed by an MySQL database.
var OutboxRepository outbox.Repository = mysqloutbox.Repository{}
