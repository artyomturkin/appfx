package appfx

import "go.uber.org/fx"

var Module = fx.Options(
	fx.NopLogger,
	configOptions,
	otelOptions,
	loggingOptions,
	routerOptions,
)
