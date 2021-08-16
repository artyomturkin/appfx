package appfx

import (
	"context"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/hashicorp/go-hclog"
	"go.uber.org/config"
	"go.uber.org/fx"
)

var loggingOptions = fx.Options(
	fx.Provide(provideHCLLogger),
	fx.Provide(watermillLogger),
	fx.NopLogger,
	fx.ErrorHook(errHandler{}),

	fx.Invoke(func(lc fx.Lifecycle, logger hclog.Logger) {
		lc.Append(fx.Hook{
			OnStart: func(_ context.Context) error {
				logger.Info("running")
				return nil
			},
			OnStop: func(_ context.Context) error {
				logger.Info("stopping")
				return nil
			},
		})
	}),
)

type logging struct {
	Level       string `yaml:"level"`
	SystemLevel string `yaml:"system_level"`
	JSON        bool   `yaml:"json"`
}

var defaultLogging = logging{
	Level:       "INFO",
	SystemLevel: "ERROR",
	JSON:        false,
}

type errHandler struct {
	fx.ErrorHandler
}

func (e errHandler) HandleError(err error) {
	hclog.Default().Error("failed to build fx app", "err", err)
}

func provideHCLLogger(c config.Provider) (hclog.Logger, error) {
	config := defaultLogging
	if err := c.Get("logging").Populate(&config); err != nil {
		return nil, err
	}

	level := hclog.Info

	switch config.Level {
	case "DEBUG":
		level = hclog.Debug
	case "ERROR":
		level = hclog.Error
	case "WARN":
		level = hclog.Warn
	case "TRACE":
		level = hclog.Trace
	}

	return hclog.New(&hclog.LoggerOptions{
		Level:      level,
		JSONFormat: config.JSON,
	}), nil
}

func watermillLogger(c config.Provider) (watermill.LoggerAdapter, error) {
	config := defaultLogging
	if err := c.Get("logging").Populate(&config); err != nil {
		return nil, err
	}

	level := hclog.Error

	switch config.SystemLevel {
	case "DEBUG":
		level = hclog.Debug
	case "ERROR":
		level = hclog.Error
	case "WARN":
		level = hclog.Warn
	case "TRACE":
		level = hclog.Trace
	}

	logger := hclog.New(&hclog.LoggerOptions{
		Level:      level,
		JSONFormat: config.JSON,
	}).Named("system")

	return &hclogAdapter{logger: logger}, nil
}

type hclogAdapter struct {
	logger hclog.Logger
}

func (h *hclogAdapter) Error(msg string, err error, fields watermill.LogFields) {
	for k, v := range fields {
		h.logger.Error(msg, k, v, "error", err)
	}
}
func (h *hclogAdapter) Info(msg string, fields watermill.LogFields) {
	for k, v := range fields {
		h.logger.Info(msg, k, v)
	}
}
func (h *hclogAdapter) Debug(msg string, fields watermill.LogFields) {
	for k, v := range fields {
		h.logger.Debug(msg, k, v)
	}
}
func (h *hclogAdapter) Trace(msg string, fields watermill.LogFields) {
	for k, v := range fields {
		h.logger.Trace(msg, k, v)
	}
}

func (h *hclogAdapter) With(fields watermill.LogFields) watermill.LoggerAdapter {
	logg := h.logger
	for k, v := range fields {
		logg = logg.With(k, v)
	}

	return &hclogAdapter{logger: logg}
}
