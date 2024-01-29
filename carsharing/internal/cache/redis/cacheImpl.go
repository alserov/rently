package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/alserov/rently/carsharing/internal/cache"
	"github.com/redis/go-redis/v9"
)

func NewCache(db *redis.Client) cache.Cache {
	return &repository{db: db}
}

type repository struct {
	db *redis.Client
}

func (r repository) Set(ctx context.Context, key string, val interface{}, exp time.Duration) error {
	b, err := json.Marshal(val)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %v", err)
	}

	v := r.db.Set(ctx, key, b, exp)
	if err = v.Err(); err != nil {
		return fmt.Errorf("failed to set to cache: %v", err)
	}
	return nil
}

func (r repository) Get(ctx context.Context, key string, target any) error {
	val, err := r.db.Get(ctx, key).Result()
	if err != nil {
		return fmt.Errorf("failed to get from cache: %v", err)
	}

	if err = json.Unmarshal([]byte(val), &target); err != nil {
		return fmt.Errorf("failed to unmarshal value: %v", err)
	}

	return nil
}
