package saga

import "github.com/jmalloc/ax/routing"

var _ routing.MessageHandler = (*MessageHandler)(nil) // ensure MessageHandler implements MessageHandler
