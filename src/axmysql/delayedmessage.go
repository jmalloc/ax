package axmysql

import (
	"github.com/jmalloc/ax/src/ax/delayedmessage"
	mysqldelayedmessage "github.com/jmalloc/ax/src/axmysql/delayedmessage"
)

// DelayedMessageRepository is an delayed message repository backed by a MySQL database.
var DelayedMessageRepository delayedmessage.Repository = mysqldelayedmessage.Repository{}
