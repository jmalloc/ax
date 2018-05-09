package logging

import (
	"github.com/jmalloc/ax/src/ax"
)

func formatMessageType(mt ax.MessageType) string {
	if mt.IsCommand() {
		return mt.Name + "?"
	} else if mt.IsEvent() {
		return mt.Name + "!"
	}

	return mt.Name
}
