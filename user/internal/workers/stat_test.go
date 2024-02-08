package workers

import (
	brokermock "github.com/alserov/rently/user/internal/utils/broker/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestStatMetricsWorker_Action(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	brkr := brokermock.NewMockProducer(ctrl)
	brkr.EXPECT().
		Produce(gomock.Any(), gomock.Any(), gomock.Eq(GOROUTINE_METRIC_ID), gomock.Any()).
		Return(nil).
		Times(1)

	w := NewStatMetricWorker(StatMetricParams{
		Producer: brkr,
	})

	t.Run("StatMetricWorkerAction", func(t *testing.T) {
		err := w.Action()
		require.NoError(t, err)
	})
}
