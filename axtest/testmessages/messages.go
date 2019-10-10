package testmessages

import "errors"

// MessageDescription returns a human-readable description of the message.
func (*Message) MessageDescription() string { return "test message" }

// MessageDescription returns a human-readable description of the message.
func (*MessageA) MessageDescription() string { return "test message A" }

// MessageDescription returns a human-readable description of the message.
func (*MessageB) MessageDescription() string { return "test message B" }

// MessageDescription returns a human-readable description of the message.
func (*MessageC) MessageDescription() string { return "test message C" }

// MessageDescription returns a human-readable description of the message.
func (*MessageD) MessageDescription() string { return "test message D" }

// MessageDescription returns a human-readable description of the message.
func (*MessageE) MessageDescription() string { return "test message E" }

// MessageDescription returns a human-readable description of the message.
func (*MessageF) MessageDescription() string { return "test message F" }

// IsCommand marks the message as a command.
func (*Command) IsCommand() {}

// MessageDescription returns a human-readable description of the message.
func (*Command) MessageDescription() string { return "test command" }

// IsCommand marks the message as a command.
func (*SelfValidatingCommand) IsCommand() {}

// Validate implements the endpoint.SelfValidatingMessage interface
func (*SelfValidatingCommand) Validate() error { return nil }

// MessageDescription returns a human-readable description of the message.
func (*SelfValidatingCommand) MessageDescription() string { return "test self-validating command" }

// IsCommand marks the message as a command.
func (*FailedSelfValidatingCommand) IsCommand() {}

// Validate implements the endpoint.SelfValidatingMessage interface.
// It always returns a test validation message.
func (*FailedSelfValidatingCommand) Validate() error {
	return errors.New("test command validation error")
}

// MessageDescription returns a human-readable description of the message.
func (*FailedSelfValidatingCommand) MessageDescription() string {
	return "test self-validating failing command"
}

// IsEvent marks the message as an event.
func (*Event) IsEvent() {}

// MessageDescription returns a human-readable description of the message.
func (*Event) MessageDescription() string { return "test event" }

// IsEvent marks the message as an event.
func (*SelfValidatingEvent) IsEvent() {}

// Validate implements the endpoint.SelfValidatingMessage interface
func (*SelfValidatingEvent) Validate() error { return nil }

// MessageDescription returns a human-readable description of the message.
func (*SelfValidatingEvent) MessageDescription() string { return "test self-validating event" }

// IsEvent marks the message as an event.
func (*FailedSelfValidatingEvent) IsEvent() {}

// Validate implements the endpoint.SelfValidatingMessage interface.
// It always returns a test validation message.
func (*FailedSelfValidatingEvent) Validate() error {
	return errors.New("test command validation error")
}

// MessageDescription returns a human-readable description of the message.
func (*FailedSelfValidatingEvent) MessageDescription() string {
	return "test self-validating failing event"
}
