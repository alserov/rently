package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"time"
)

type Params struct {
	Addr     string
	Password string
}

func MustConnect(p Params) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     p.Addr,
		Password: p.Password, // no password set
		DB:       0,          // use default DB
	})

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		panic("failed to ping cache connection: " + err.Error())
	}

	return client
}
