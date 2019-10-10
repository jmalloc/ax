package projection_test

import (
	. "github.com/jmalloc/ax/projection"
	"github.com/jmalloc/ax/routing"
)

var _ routing.MessageHandler = (*MessageHandler)(nil) // ensure MessageHandler implements MessageHandler
