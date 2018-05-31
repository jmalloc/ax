package aggregate

// Option is a function that applies some change to the underlying saga
// implementation of an aggregate.
type Option func(*Saga)

// IdentifyByField returns an aggregate option that maps commands to instances
// by using the value of the message field named n as the aggregate ID.
func IdentifyByField(n string) Option {
	return func(sg *Saga) {
		sg.Identifier = &ByFieldIdentifier{
			FieldName: n,
		}
	}
}
