package main

import (
	"time"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/hashicorp/go-hclog"
	"go.uber.org/fx"
)

type termHandlerConfig struct {
	fx.In

	Log hclog.Logger
	R   *message.Router

	Sub message.Subscriber
	Pub message.Publisher
}

func newTermHandler(cfg inputHandlerConfig) {
	name := "term"
	log := cfg.Log.Named(name)

	cfg.R.AddNoPublisherHandler(
		name,
		"last",
		cfg.Sub,
		last(log),
	)
}

func last(log hclog.Logger) func(*message.Message) error {
	return func(msg *message.Message) error {
		log.Info("got message")

		time.Sleep(200 * time.Millisecond)

		log.Info("processed message")

		return nil
	}
}
