package ax

import "time"

// Delay is an option that delays sending the message until a duration has
// passed. Events can not be delayed.
func Delay(d time.Duration) ExecuteOption {
	t := time.Now().Add(d)
	return DelayUntil(t)
}

// DelayUntil is an option that delays sending the message until a specific
// time. Events can not be delayed.
func DelayUntil(t time.Time) ExecuteOption {
	return delayOption{t}
}

// delayOption provides the implementation of SendOption for the Delay and
// DelayUntil options.
type delayOption struct {
	Time time.Time
}

func (o delayOption) ApplyExecuteOption(env *Envelope) error {
	env.SendAt = o.Time
	return nil
}
