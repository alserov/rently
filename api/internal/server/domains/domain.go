package domains

import "time"

type Params[T any] struct {
	Client       T
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}
