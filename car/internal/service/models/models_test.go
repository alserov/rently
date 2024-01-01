package models

import (
	"fmt"
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
	fmt.Println(period)
	require.NotEmpty(t, period)
}
