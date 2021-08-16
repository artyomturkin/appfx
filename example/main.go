package main

import (
	"github.com/artyomturkin/appfx"
	"go.uber.org/fx"
)

func main() {

	fx.New(
		appfx.Module,

		fx.Invoke(newInputHandler),
		fx.Invoke(newIntermidHandler),
		fx.Invoke(newTermHandler),
	).Run()
}
