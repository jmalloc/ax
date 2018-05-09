package observability_test

import (
	. "github.com/jmalloc/ax/src/ax/observability"
)

var (
	ensureLoggingObserverIsBeforeInboundObserver BeforeInboundObserver = &LoggingObserver{}
	ensureLoggingObserverIsAfterInboundObserver  AfterInboundObserver  = &LoggingObserver{}
	ensureLoggingObserverIsAfterOutboundObserver AfterOutboundObserver = &LoggingObserver{}
)
