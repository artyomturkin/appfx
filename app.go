package appfx

import (
	"context"

	"github.com/go-playground/validator/v10"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.uber.org/config"
	"go.uber.org/fx"
)

var Module = fx.Options(
	// fx.NopLogger,
	configOptions,
	otelOptions,
	loggingOptions,
	routerOptions,
	fx.Provide(appToResource),
)

type Application struct {
	Name        string `validate:"required"`
	Version     string `validate:"required"`
	Namespace   string `validate:"required"`
	Environment string `validate:"required"`
}

func appToResource(config config.Provider) (*resource.Resource, error) {
	var app Application
	if err := config.Get("application").Populate(&app); err != nil {
		return nil, err
	}

	err := validator.New().Struct(app)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()

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
