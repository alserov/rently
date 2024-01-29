package cache

import (
	"context"
	"time"
)

type Cache interface {
	Set(ctx context.Context, key string, val interface{}, exp time.Duration) error
	Get(ctx context.Context, key string, target any) error
}
