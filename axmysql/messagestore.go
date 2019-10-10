package axmysql

import (
	mysqlmessagestore "github.com/jmalloc/ax/axmysql/messagestore"
	"github.com/jmalloc/ax/messagestore"
)

// MessageStore is a message store backed by a MySQL database.
var MessageStore messagestore.GloballyOrderedStore = mysqlmessagestore.Store{}
