package axcli

import (
	"context"
	"fmt"
	"time"

	"github.com/jmalloc/ax/src/ax"
	"github.com/spf13/cobra"
	strcase "github.com/stoewer/go-strcase"
)

// NewCommands returns  new CLI commands for each of the given message types.
func NewCommands(
	s ax.Sender,
	t ax.MessageTypeSet,
) ([]*cobra.Command, error) {
	var cmds []*cobra.Command

	for _, mt := range t.Members() {
		cmd, err := NewCommand(s, mt)
		if err != nil {
			return nil, err
		}

		cmds = append(cmds, cmd)
	}

	return cmds, nil
}

// NewCommand generates a CLI command for given message type.
func NewCommand(s ax.Sender, mt ax.MessageType) (*cobra.Command, error) {
	m := mt.New()
	var usage string
	var send func(ctx context.Context) error

	switch v := m.(type) {
	case ax.Command:
		usage = fmt.Sprintf(
			"Execute the '%s' command",
			mt.Name,
		)
		send = func(ctx context.Context) error {
			return s.ExecuteCommand(ctx, v)
		}
	case ax.Event:
		usage = fmt.Sprintf(
			"Publish the '%s' event",
			mt.Name,
		)
		send = func(ctx context.Context) error {
			return s.PublishEvent(ctx, v)
		}
	default:
		return nil, fmt.Errorf(
			"%s is neither a command nor an event",
			mt.Name,
		)
	}

	cmd := &cobra.Command{
		Use:   strcase.KebabCase(mt.MessageName()),
		Short: usage,
		Args:  cobra.NoArgs,
		RunE: func(c *cobra.Command, args []string) error {
			timeout, err := c.Flags().GetDuration("timeout")
			if err != nil {
				return err
			}
			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()

			c.SilenceUsage = true
			return send(ctx)
		},
	}

	flags := cmd.Flags()

	flags.DurationP(
		"timeout", "t",
		5*time.Second,
		"sets the timeout for command execution",
	)

	if err := declareFlags(flags, m); err != nil {
		return nil, err
	}

	return cmd, nil
}
