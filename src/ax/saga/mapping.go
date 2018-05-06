package saga

import "github.com/jmalloc/ax/src/ax"

// MappingKey is an identifier that is used to map incoming messages to
// saga instances.
//
// For each incoming message m, Saga.MapMessage(m) is called to produced a
// mapping key k. A 2-tuple containing m's message type and k is used as a
// identifier to find a saga instance with the same message type / key
// combination.
//
// If there is no matching saga instance then either a new instance is created
// (if m's message type is marked as a trigger), or the saga's not-found handler
// is called.
//
// Whenever a saga instance is modified, a MappingTable is produced by calling
// Saga.MapInstance() for each of the saga's message types. There can be only
// one saga instance for each message type / key combination.
type MappingKey string

// MappingTable is a map of message type to mapping key for a specific saga
// instance.
type MappingTable map[ax.MessageType]MappingKey

// buildMappingTable returns a mapping table for i containing each of its
// supported message types.
func buildMappingTable(s Saga, i Instance) MappingTable {
	triggers, others := s.MessageTypes()
	types := triggers.Union(others)
	table := MappingTable{}

	for _, mt := range types.Members() {
		table[mt] = s.MapInstance(mt, i)
	}

	return table
}
