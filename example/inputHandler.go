package main

import (
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/hashicorp/go-hclog"
	"go.uber.org/fx"
)

type inputHandlerConfig struct {
	fx.In

	Log hclog.Logger
	R   *message.Router

	Sub message.Subscriber
	Pub message.Publisher
}

func newInputHandler(cfg inputHandlerConfig) {
	name := "input"
	log := cfg.Log.Named(name)

	cfg.R.AddHandler(
		name,
		"input",
		cfg.Sub,
		"inter",
		cfg.Pub,
		passthrough(log),
	)
}
