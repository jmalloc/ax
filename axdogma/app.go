package axdogma

import (
	"context"

	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/enginekit/config"
	"github.com/jmalloc/ax/projection"
	"github.com/jmalloc/ax/routing"
	"github.com/jmalloc/ax/saga"
)

// App contains Dogma application components adapted to Ax interfaces.
type App struct {
	Aggregates   []saga.Saga
	Processes    []saga.Saga
	Integrations []routing.MessageHandler
	Projections  []projection.Projector
}

// New returns a structure that contains a Dogma application's aggregates,
// processes, integrations and projections adapted into the most appropriate Ax
// type.
func New(app dogma.Application) (*App, error) {
	cfg, err := config.NewApplicationConfig(app)
	if err != nil {
		return nil, err
	}

	v := &visitor{}
	err = cfg.Accept(context.Background(), v)
	return &v.app, err
}

type visitor struct {
	app App
}

func (v *visitor) VisitApplicationConfig(ctx context.Context, cfg *config.ApplicationConfig) error {
	for _, hcfg := range cfg.HandlersByName {
		if err := hcfg.Accept(ctx, v); err != nil {
			return err
		}
	}

	return nil
}

func (v *visitor) VisitAggregateConfig(_ context.Context, cfg *config.AggregateConfig) error {
	a := &AggregateAdaptor{
		Key:     cfg.HandlerIdentity.Key,
		Handler: cfg.Handler,
	}

	for mt := range cfg.ConsumedMessageTypes() {
		a.CommandTypes = a.CommandTypes.Add(
			convertMessageType(mt),
		)
	}

	v.app.Aggregates = append(v.app.Aggregates, a)

	return nil
}

func (v *visitor) VisitProcessConfig(_ context.Context, cfg *config.ProcessConfig) error {
	a := &ProcessAdaptor{
		Key:     cfg.HandlerIdentity.Key,
		Handler: cfg.Handler,
	}

	for mt := range cfg.ConsumedMessageTypes() {
		a.EventTypes = a.EventTypes.Add(
			convertMessageType(mt),
		)
	}

	v.app.Processes = append(v.app.Processes, a)

	return nil
}

func (v *visitor) VisitIntegrationConfig(_ context.Context, cfg *config.IntegrationConfig) error {
	a := &IntegrationAdaptor{
		Handler: cfg.Handler,
	}

	for mt := range cfg.ConsumedMessageTypes() {
		a.CommandTypes = a.CommandTypes.Add(
			convertMessageType(mt),
		)
	}

	v.app.Integrations = append(v.app.Integrations, a)

	return nil
}

func (v *visitor) VisitProjectionConfig(_ context.Context, cfg *config.ProjectionConfig) error {
	a := &ProjectionAdaptor{
		Key:     cfg.HandlerIdentity.Key,
		Handler: cfg.Handler,
	}

	for mt := range cfg.ConsumedMessageTypes() {
		a.EventTypes = a.EventTypes.Add(
			convertMessageType(mt),
		)
	}

	v.app.Projections = append(v.app.Projections, a)

	return nil
}
