package messagetest

// Description returns a human-readable description of the message.
func (*Message) Description() string { return "test message" }

// IsCommand marks the message as a command.
func (*Command) IsCommand() {}

// Description returns a human-readable description of the message.
func (*Command) Description() string { return "test command" }

// IsEvent marks the message as an event.
func (*Event) IsEvent() {}

// Description returns a human-readable description of the message.
func (*Event) Description() string { return "test event" }
