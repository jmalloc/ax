package messagetest

// MessageDescription returns a human-readable description of the message.
func (*Message) MessageDescription() string { return "test message" }

// IsCommand marks the message as a command.
func (*Command) IsCommand() {}

// MessageDescription returns a human-readable description of the message.
func (*Command) MessageDescription() string { return "test command" }

// IsEvent marks the message as an event.
func (*Event) IsEvent() {}

// MessageDescription returns a human-readable description of the message.
func (*Event) MessageDescription() string { return "test event" }
