package saga_test

import (
	"github.com/jmalloc/ax/src/ax/routing"
	. "github.com/jmalloc/ax/src/ax/saga"
)

var ensureHandlerIsMessageHandler routing.MessageHandler = &MessageHandler{}
