package axdogma

import (
	"context"

	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/enginekit/config"
	"github.com/jmalloc/ax/src/ax/projection"
	"github.com/jmalloc/ax/src/ax/routing"
)

// FromApp builds a set of Ax message handlers and projectors for the given
// Dogma application.
func FromApp(
	app dogma.Application,
) (
	[]routing.MessageHandler,
	[]projection.Projector,
	error,
) {
	cfg, err := config.NewApplicationConfig(app)
	if err != nil {
		return nil, nil, err
	}

	v := &visitor{}
	err = cfg.Accept(context.Background(), v)

	return v.handlers, v.projectors, err
}

type visitor struct {
	handlers   []routing.MessageHandler
	projectors []projection.Projector
}

func (v *visitor) VisitApplicationConfig(ctx context.Context, cfg *config.ApplicationConfig) error {
	for _, hcfg := range cfg.Handlers {
		if err := hcfg.Accept(ctx, v); err != nil {
			return err
		}
	}

	return nil
}

func (v *visitor) VisitAggregateConfig(_ context.Context, cfg *config.AggregateConfig) error {
	panic("not implemented")
}

func (v *visitor) VisitProcessConfig(_ context.Context, cfg *config.ProcessConfig) error {
	panic("not implemented")
}

func (v *visitor) VisitIntegrationConfig(_ context.Context, cfg *config.IntegrationConfig) error {
	panic("not implemented")
}

func (v *visitor) VisitProjectionConfig(_ context.Context, cfg *config.ProjectionConfig) error {
	a := &ProjectionAdaptor{
		Name:    cfg.HandlerName,
		Handler: cfg.Handler,
	}

	for mt := range cfg.EventTypes() {
		a.EventTypes.Add(
			convertMessageType(mt),
		)
	}

	v.projectors = append(v.projectors, a)

	return nil
}
