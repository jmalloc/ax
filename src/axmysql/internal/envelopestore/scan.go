package envelopestore

import (
	"github.com/jmalloc/ax/src/ax"
	"github.com/jmalloc/ax/src/ax/endpoint"
	"github.com/jmalloc/ax/src/ax/marshaling"
)

// Scanner is an interface for scanning values from a row or rows.
type Scanner interface {
	Scan(v ...interface{}) error
}

// Columns is the ordered set of columns that must be SELECTed for use with Scan().
const Columns = `message_id,
				causation_id,
				correlation_id,
				created_at,
				send_at,
				content_type,
				data,
				operation,
				destination`

// Scan constructs an outbound envelope by scanning values using s.
func Scan(s Scanner) (endpoint.OutboundEnvelope, error) {
	var env endpoint.OutboundEnvelope

	var (
		ct        string
		data      []byte
		createdAt string
		sendAt    string
	)

	err := s.Scan(
		&env.MessageID,
		&env.CausationID,
		&env.CorrelationID,
		&createdAt,
		&sendAt,
		&ct,
		&data,
		&env.Operation,
		&env.DestinationEndpoint,
	)
	if err != nil {
		return endpoint.OutboundEnvelope{}, err
	}

	err = marshaling.UnmarshalTime(createdAt, &env.CreatedAt)
	if err != nil {
		return endpoint.OutboundEnvelope{}, err
	}

	err = marshaling.UnmarshalTime(sendAt, &env.SendAt)
	if err != nil {
		return endpoint.OutboundEnvelope{}, err
	}

	env.Message, err = ax.UnmarshalMessage(ct, data)

	return env, err
}
