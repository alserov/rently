package metrics

import (
	"context"
)

type Worker interface {
	MustStart(ctx context.Context)
}
