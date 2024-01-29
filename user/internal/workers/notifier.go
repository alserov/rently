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

func NewNotifier(ticker *time.Ticker, actor Actor) Worker {
	return &notifier{
		t:     ticker,
		actor: actor,
	}
}

type notifier struct {
	t     *time.Ticker
	actor Actor
}

func (n *notifier) MustStart() {
	for range n.t.C {
		n.actor.Notify()
	}
}
