package appfx

import "go.uber.org/fx"

var Module = fx.Options(
	fx.NopLogger,
	configOptions,
	otelOptions,
	loggingOptions,
	routerOptions,
)

type Application struct {
	Name        string `yaml:"name"        validate:"required"`
	Version     string `yaml:"version"     validate:"required"`
	Namespace   string `yaml:"namespace"   validate:"required"`
	Environment string `yaml:"environment" validate:"required"`
}
