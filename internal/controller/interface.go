package controller

import (
	"context"
)

type Controller interface {
	ShouldBeRunning() bool
	Start(ctx context.Context) error
	Stop(_ context.Context) error
}
