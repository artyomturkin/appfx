package main

import (
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/hashicorp/go-hclog"
	"go.uber.org/fx"
)

type intermidHandlerConfig struct {
	fx.In

	Log hclog.Logger
	R   *message.Router

	Sub message.Subscriber
	Pub message.Publisher
}

func newIntermidHandler(cfg inputHandlerConfig) {
	name := "intermid"
	log := cfg.Log.Named(name)

	cfg.R.AddHandler(
		name,
		"inter",
		cfg.Sub,
		"last",
		cfg.Pub,
		passthrough(log),
	)
}
