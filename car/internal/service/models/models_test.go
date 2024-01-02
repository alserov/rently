package models

import (
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestCreateRentReq_Period(t *testing.T) {
	start := time.Now()
	end := start.Add(time.Hour * 24 * 3)

	rent := CreateRentReq{
		RentStart: &start,
		RentEnd:   &end,
	}

	period := rent.Period()
	require.Equal(t, time.Hour*24*3, period)
}
