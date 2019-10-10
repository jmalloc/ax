package routing

import (
	"fmt"

	"github.com/jmalloc/ax"
)

// HandlerTable is a set of rules that determines which handlers receive a
// message of a specific type.
type HandlerTable map[ax.MessageType][]MessageHandler

// NewHandlerTable returns a handler table that locates message handlers
// based on the message types that they handle.
func NewHandlerTable(handlers ...MessageHandler) (HandlerTable, error) {
	ht := HandlerTable{}

	for _, h := range handlers {
		for _, mt := range h.MessageTypes().Members() {
			x := ht[mt]

			if mt.IsCommand() && len(x) != 0 {
				return nil, fmt.Errorf(
					"can not build handler table, multiple message handlers are defined for the '%s' command",
					mt.Name,
				)
			}

			ht[mt] = append(x, h)
		}
	}

	return ht, nil
}

// Lookup returns the message handlers that handle mt.
func (ht HandlerTable) Lookup(mt ax.MessageType) []MessageHandler {
	return ht[mt]
}
