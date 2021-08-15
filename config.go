package appfx

import (
	"context"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/artyomturkin/go-from-uri/kafka"
	"github.com/sherifabdlnaby/configuro"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.uber.org/fx"
)

type config struct {
	Application   Application   `yaml:"application"   validate:"required"`
	Observability Observability `yaml:"observability" validate:"required"`

	Subscriber string `yaml:"subscriber"`
	Publisher  string `yaml:"publisher"`
}

func parseConfig() (config, error) {
	var conf config

	co := []configuro.ConfigOptions{}

	c, err := configuro.NewConfig(co...)
	if err != nil {
		return conf, err
	}

	err = c.Load(&conf)
	if err != nil {
		return conf, err
	}

	return conf, nil
}

func appToResource(ctx context.Context, app Application) (*resource.Resource, error) {
	return resource.New(ctx,
		resource.WithFromEnv(),
		resource.WithProcess(),
		resource.WithTelemetrySDK(),
		resource.WithHost(),
		resource.WithAttributes(
			semconv.ServiceNameKey.String(app.Name),
			semconv.ServiceVersionKey.String(app.Version),
			semconv.DeploymentEnvironmentKey.String(app.Environment),
			semconv.ServiceNamespaceKey.String(app.Namespace),
		),
	)
}

func configure() fx.Option {
	cfg, err := parseConfig()
	if err != nil {
		return fx.Error(err)
	}

	ctx := context.Background()
	opts := []fx.Option{}

	r, err := appToResource(ctx, cfg.Application)
	if err != nil {
		return fx.Error(err)
	}

	opts = append(opts, fx.Provide(func(lc fx.Lifecycle, logger watermill.LoggerAdapter) (message.Publisher, error) {
		res, err := kafka.NewWatermillPublisher(cfg.Publisher, logger)
		if err != nil {
			return nil, err
		}

		lc.Append(fx.Hook{
			OnStop: func(_ context.Context) error {
				return res.Close()
			},
		})

		return res, nil
	}))

	opts = append(opts, fx.Provide(func(lc fx.Lifecycle, logger watermill.LoggerAdapter) (message.Subscriber, error) {
		res, err := kafka.NewWatermillSubscriber(cfg.Subscriber, logger)
		if err != nil {
			return nil, err
		}

		lc.Append(fx.Hook{
			OnStop: func(_ context.Context) error {
				return res.Close()
			},
		})

		return res, nil
	}))

	opts = append(
		opts,
		fx.Supply(r),
		fx.Supply(cfg.Observability.Logging),
		fx.Supply(cfg.Observability.Metrics),
		fx.Supply(cfg.Observability.Tracing),
	)

	return fx.Options(opts...)
}

var configOptions = configure()
