package axdogma

import fmt "fmt"

// IsEvent marks the message as an event.
func (m *Timeout) IsEvent() {}

// MessageDescription returns a human-readable description of the message.
func (m *Timeout) MessageDescription() string {
	return fmt.Sprintf(
		"timeout for %s %s",
		m.Process,
		m.ProcessId,
	)
}
