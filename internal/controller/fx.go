package controller

import (
	"go.uber.org/fx"
)

func RunControllerFx(lc fx.Lifecycle, ctrl Controller) {
	if ctrl.ShouldBeRunning() {
		lc.Append(fx.Hook{
			OnStart: ctrl.Start,
			OnStop:  ctrl.Stop,
		})
	}
}

func AsController(f any) any {
	return fx.Annotate(
		f,
		fx.As(new(Controller)),
		fx.ResultTags(`group:"controllers"`),
	)
}

type RunControllersFxParams struct {
	fx.In
	Controllers []Controller `group:"controllers"`
	Lifecycle   fx.Lifecycle
}

func RunControllersFx(params RunControllersFxParams) {
	for _, ctrl := range params.Controllers {
		RunControllerFx(params.Lifecycle, ctrl)
	}
}
