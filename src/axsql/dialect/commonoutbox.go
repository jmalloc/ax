package dialect

import (
	"database/sql"

	"github.com/jmalloc/ax/src/ax"
	"github.com/jmalloc/ax/src/ax/bus"
	"github.com/jmalloc/ax/src/ax/marshaling"
)

func scanOutboxMessages(
	rows *sql.Rows,
	causationID ax.MessageID,
) ([]bus.OutboundEnvelope, error) {
	var envelopes []bus.OutboundEnvelope

	for rows.Next() {
		env := bus.OutboundEnvelope{
			Envelope: ax.Envelope{
				CausationID: causationID,
			},
		}

		if err := scanOutboxMessage(rows, &env); err != nil {
			return nil, err
		}

		envelopes = append(envelopes, env)
	}

	return envelopes, nil
}

func scanOutboxMessage(rows *sql.Rows, env *bus.OutboundEnvelope) error {
	var (
		ct   string
		body []byte
	)

	err := rows.Scan(
		&env.MessageID,
		&env.CorrelationID,
		&env.Time,
		&ct,
		&body,
		&env.Operation,
		&env.DestinationEndpoint,
	)
	if err != nil {
		return err
	}

	env.Message, err = marshaling.UnmarshalMessage(ct, body)

	return err
}
