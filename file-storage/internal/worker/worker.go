package worker

import (
	"context"
)

type Worker interface {
	MustStart(ctx context.Context)
}
