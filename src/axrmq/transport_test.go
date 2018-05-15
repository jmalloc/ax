package axrmq_test

import (
	"github.com/jmalloc/ax/src/ax/endpoint"
	. "github.com/jmalloc/ax/src/axrmq"
)

var ensureTransportIsEndpointTransport endpoint.Transport = &Transport{}
