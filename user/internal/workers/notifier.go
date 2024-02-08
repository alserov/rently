package workers

import (
	"fmt"
	"github.com/alserov/rently/user/internal/log"
	"time"
)

type Worker interface {
	MustStart()
}

type Actor interface {
	Action() error
}

func StartWithTicker(ticker *time.Ticker, actor Actor) {
	l := log.GetLogger()
	for range ticker.C {
		if err := actor.Action(); err != nil {
			l.Error(fmt.Errorf("ticker error: %w", err).Error())
		}
	}
}
