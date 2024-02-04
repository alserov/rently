package workers

import (
	"time"
)

type Worker interface {
	MustStart()
}

type Actor interface {
	Notify()
}

func StartNotifier(ticker *time.Ticker, actor Actor) {
	for range ticker.C {
		actor.Notify()
	}
}
