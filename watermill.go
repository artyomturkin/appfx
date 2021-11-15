package appfx

import (
	"context"
	"fmt"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	watermilluri "github.com/artyomturkin/go-from-uri/watermill"
	watermillext "github.com/artyomturkin/watermill-extension"
	"go.uber.org/config"
	"go.uber.org/fx"
)

var routerOptions = fx.Options(
	fx.Provide(router),
	fx.Provide(buildSubscriber),
	fx.Provide(buildPublisher),
)

func router(lc fx.Lifecycle, logger watermill.LoggerAdapter) (*message.Router, error) {
	router, err := message.NewRouter(message.RouterConfig{}, watermill.NopLogger{})
	if err != nil {
		return nil, err
	}

	router.AddMiddleware(
		watermillext.OpenTelemetryMiddleware,
		middleware.Recoverer,
	)

	rctx, cancel := context.WithCancel(context.Background())

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() { router.Run(rctx) }()
			return nil
		},
		OnStop: func(_ context.Context) error {
			cancel()
			return router.Close()
		},
	})

	return router, nil
}

func buildPublisher(lc fx.Lifecycle, logger watermill.LoggerAdapter, c config.Provider) (message.Publisher, error) {
	var str string
	if err := c.Get("publisher").Populate(&str); err == nil && str != "" {
		res, err := watermilluri.NewPublisher(str, logger)
		if err != nil {
			return nil, err
		}

		lc.Append(fx.Hook{
			OnStop: func(_ context.Context) error {
				return res.Close()
			},
		})

		return res, nil
	}
	return nil, fmt.Errorf("subscriber config not set")
}

func buildSubscriber(lc fx.Lifecycle, logger watermill.LoggerAdapter, c config.Provider) (message.Subscriber, error) {
	var str string
	if err := c.Get("subscriber").Populate(&str); err == nil && str != "" {
		res, err := watermilluri.NewSubscriber(str, logger)
		if err != nil {
			return nil, err
		}

		lc.Append(fx.Hook{
			OnStop: func(_ context.Context) error {
				return res.Close()
			},
		})

		return res, nil
	}

	return nil, fmt.Errorf("subscriber config not set")
}
