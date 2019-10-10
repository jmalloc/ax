package marshaling

import "time"

// MarshalTime marshals t to an RFC3339+Nano string.
func MarshalTime(t time.Time) string {
	return t.Format(time.RFC3339Nano)
}

// UnmarshalTime unmarshals an RFC3339+Nano formatted string into t.
func UnmarshalTime(s string, t *time.Time) error {
	v, err := time.Parse(time.RFC3339Nano, s)
	if err != nil {
		return err
	}

	*t = v
	return nil
}
