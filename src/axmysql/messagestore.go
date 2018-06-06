package axmysql

import (
	"github.com/jmalloc/ax/src/ax/messagestore"
	mysqlmessagestore "github.com/jmalloc/ax/src/axmysql/messagestore"
)

// MessageStore is a message store backed by a MySQL database.
var MessageStore messagestore.GloballyOrderedStore = mysqlmessagestore.Store{}
