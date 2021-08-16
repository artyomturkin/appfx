package main

import (
	"time"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/hashicorp/go-hclog"
)

func passthrough(log hclog.Logger) func(msg *message.Message) ([]*message.Message, error) {
	return func(msg *message.Message) ([]*message.Message, error) {
		log.Info("got message")

		time.Sleep(200 * time.Millisecond)

		log.Info("processed message")

		return []*message.Message{msg}, nil
	}
}
