package axrmq_test

import (
	. "github.com/jmalloc/ax/axrmq"
	"github.com/jmalloc/ax/endpoint"
)

var (
	_ endpoint.InboundTransport  = (*Transport)(nil) // ensure Transport implements InboundTransport
	_ endpoint.OutboundTransport = (*Transport)(nil) // ensure Transport implements OutboundTransport
)
