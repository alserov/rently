package workers

import (
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

type testActor struct {
	counter int
}

func (t *testActor) Action() error {
	t.counter++
	return nil
}

func TestStartNotifier(t *testing.T) {
	ta := testActor{}
	go StartWithTicker(time.NewTicker(time.Millisecond*100), &ta)

	start := time.Now()
	time.Sleep(time.Second + time.Millisecond*30)
	require.Equal(t, int64(ta.counter), time.Since(start).Milliseconds()/100)
}
