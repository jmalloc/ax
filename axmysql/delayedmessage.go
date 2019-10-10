package axmysql

import (
	mysqldelayedmessage "github.com/jmalloc/ax/axmysql/delayedmessage"
	"github.com/jmalloc/ax/delayedmessage"
)

// DelayedMessageRepository is an delayed message repository backed by a MySQL database.
var DelayedMessageRepository delayedmessage.Repository = mysqldelayedmessage.Repository{}
