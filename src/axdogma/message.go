package axdogma

import (
	"reflect"

	"github.com/dogmatiq/enginekit/message"
	"github.com/jmalloc/ax/src/ax"
)

func convertMessageType(mt message.Type) ax.MessageType {
	m := reflect.Zero(mt.ReflectType()).Interface()
	return ax.TypeOf(m.(ax.Message))
}
