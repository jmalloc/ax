package axmysql

import (
	mysqloutbox "github.com/jmalloc/ax/axmysql/outbox"
	"github.com/jmalloc/ax/outbox"
)

// OutboxRepository is an outbox repository backed by a MySQL database.
var OutboxRepository outbox.Repository = mysqloutbox.Repository{}
