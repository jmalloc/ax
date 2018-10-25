package axrmq_test

import (
	"github.com/jmalloc/ax/src/ax/endpoint"
	. "github.com/jmalloc/ax/src/axrmq"
)

var (
	_ endpoint.InboundTransport  = (*Transport)(nil) // ensure Transport implements InboundTransport
	_ endpoint.OutboundTransport = (*Transport)(nil) // ensure Transport implements OutboundTransport
)
