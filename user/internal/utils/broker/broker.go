package broker

import "context"

type Producer interface {
	Produce(ctx context.Context, value any, id string, q string) error
}
