package appfx

import (
	"context"
	"fmt"
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/export/metric/aggregation"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/histogram"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/config"
	"go.uber.org/fx"

	controller "go.opentelemetry.io/otel/sdk/metric/controller/basic"
	processor "go.opentelemetry.io/otel/sdk/metric/processor/basic"
	selector "go.opentelemetry.io/otel/sdk/metric/selector/simple"
)

var otelOptions = fx.Options(
	fx.Invoke(tracerProviderJaeger),
	fx.Invoke(prometheusExporter),
)

type tracing struct {
	Jaeger string `yaml:"jaeger"`
}

func tracerProviderJaeger(lc fx.Lifecycle, resource *resource.Resource, c config.Provider) error {
	var cfg tracing
	if err := c.Get("tracing").Populate(&cfg); err != nil {
		return err
	}

	// Create the Jaeger exporter
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(cfg.Jaeger)))
	if err != nil {
		return err
	}

	tp := tracesdk.NewTracerProvider(
		// Always be sure to batch in production.
		tracesdk.WithBatcher(exp),
		// Record information about this application in an Resource.
		tracesdk.WithResource(resource),
	)

	lc.Append(
		fx.Hook{
			OnStop: tp.Shutdown,
		},
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	return nil
}

type metrics struct {
	Prometheus int `yaml:"prometheus" validate:"required"`
}

func prometheusExporter(lc fx.Lifecycle, resource *resource.Resource, cp config.Provider) error {
	var cfg metrics
	if err := cp.Get("metrics").Populate(&cfg); err != nil {
		return err
	}

	config := prometheus.Config{}
	c := controller.New(
		processor.NewFactory(
			selector.NewWithHistogramDistribution(
				histogram.WithExplicitBoundaries([]float64{
					// Fast operation
					10, 25, 50, 100, 300, 1_000, 5_000, 15_000, 60_000,
					// Long operations 2, 5, 15 minutest and 1, 5, 10 hours
					120_000, 300_000, 900_000, 3600_000, 18_000_000, 180_000_000,
				}),
			),
			aggregation.CumulativeTemporalitySelector(),
			processor.WithMemory(true),
		),
		controller.WithResource(resource),
	)
	exporter, err := prometheus.New(config, c)
	if err != nil {
		return err
	}
	global.SetMeterProvider(exporter.MeterProvider())

	mux := http.NewServeMux()
	mux.HandleFunc("/metrics", exporter.ServeHTTP)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Prometheus),
		Handler: mux,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				srv.ListenAndServe()
			}()
			return nil
		},
		OnStop: srv.Shutdown,
	})

	return nil
}
