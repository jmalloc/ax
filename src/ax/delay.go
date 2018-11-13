package ax

import "time"

// Delay is an option that delays sending the message until a duration has
// passed. Events can not be delayed.
func Delay(d time.Duration) ExecuteOption {
	return delayOption{d}
}

// DelayUntil is an option that delays sending the message until a specific
// time. Events can not be delayed.
func DelayUntil(t time.Time) ExecuteOption {
	return delayUntilOption{t}
}

// delayOption provides the implementation of SendOption for the Delay option.
type delayOption struct {
	Delay time.Duration
}

func (o delayOption) ApplyExecuteOption(env *Envelope) error {
	env.SendAt = env.CreatedAt.Add(o.Delay)
	return nil
}

// delayUntilOption provides the implementation of SendOption for the DelayUntil
// options.
type delayUntilOption struct {
	Time time.Time
}

func (o delayUntilOption) ApplyExecuteOption(env *Envelope) error {
	env.SendAt = o.Time
	return nil
}
