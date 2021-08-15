package appfx

import (
	"context"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"go.uber.org/fx"
)

var routerOptions = fx.Options(
	fx.Provide(router),
)

func router(lc fx.Lifecycle, logger watermill.LoggerAdapter) (*message.Router, error) {
	router, err := message.NewRouter(message.RouterConfig{}, watermill.NopLogger{})
	if err != nil {
		return nil, err
	}

	router.AddMiddleware(
		middleware.Recoverer,
	)

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() { router.Run(ctx) }()
			return nil
		},
		OnStop: func(_ context.Context) error {
			return router.Close()
		},
	})

	return router, nil
}
