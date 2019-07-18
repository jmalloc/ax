package axdogma

import (
	"reflect"

	"github.com/dogmatiq/enginekit/message"
	"github.com/jmalloc/ax/src/ax"
)

// convertMessageType returns the equivalent Ax message type for the given Dogma
// enginekit message type.
func convertMessageType(mt message.Type) ax.MessageType {
	m := reflect.Zero(mt.ReflectType()).Interface()
	return ax.TypeOf(m.(ax.Message))
}
