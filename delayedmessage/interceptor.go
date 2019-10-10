package delayedmessage

import (
	"context"
	"time"

	"github.com/jmalloc/ax/endpoint"
	"github.com/jmalloc/ax/internal/tracing"
	"github.com/jmalloc/ax/persistence"
)

// Interceptor is an outbound pipeline stage that intercepts messages that are
// not ready to be sent.
type Interceptor struct {
	Repository Repository
	Next       endpoint.OutboundPipeline
}

// Initialize is called during initialization of the endpoint, after the
// transport is initialized. It can be used to inspect or further configure the
// endpoint as per the needs of the pipeline.
func (i *Interceptor) Initialize(ctx context.Context, ep *endpoint.Endpoint) error {
	return i.Next.Initialize(ctx, ep)
}

// Accept passes env to the next pipeline stage only if it is ready to send now,
// otherwise it stores it to be sent in the future.
func (i *Interceptor) Accept(ctx context.Context, env endpoint.OutboundEnvelope) error {
	// send immediately if sendAt <= now
	if !env.SendAt.After(time.Now()) {
		return i.Next.Accept(ctx, env)
	}

	tracing.LogEvent(
		ctx,
		"delay",
		"intercepting the message to be sent after a delay",
		tracing.Duration("delay_for", env.Delay()),
		tracing.Time("delay_until", env.SendAt),
		tracing.TypeName("pipeline_stage", i),
	)

	tx, com, err := persistence.GetOrBeginTx(ctx)
	if err != nil {
		return err
	}
	defer com.Rollback()

	if err := i.Repository.SaveMessage(ctx, tx, env); err != nil {
		return err
	}

	return com.Commit()
}
