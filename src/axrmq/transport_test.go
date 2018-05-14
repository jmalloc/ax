package axrmq_test

import (
	"github.com/jmalloc/ax/src/ax/bus"
	. "github.com/jmalloc/ax/src/axrmq"
)

var ensureTransportIsBusTransport bus.Transport = &Transport{}
