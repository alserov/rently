package payment

import (
	"fmt"
	"github.com/alserov/rently/car/internal/service/models"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

const api_key = "sk_test_51OU56CDOnc0MdcTNBwddO2cn8NrEebjfuAGjBjj9xSyKmiUO4ajJ1vZ0yBoOsAMq0HjHqCmis2niwoj2EZYCDLOA00lcCUlWxh"

func TestPayer(t *testing.T) {
	p := NewPayer(api_key)

	const (
		amount = 100_00
		source = "tok_visa"
	)

	// invalid price
	chargeID, err := p.Debit(source, 0)
	require.Error(t, err)
	require.Empty(t, chargeID)

	// valid price
	chargeID, err = p.Debit(source, amount)
	require.NoError(t, err)
	require.NotEmpty(t, chargeID)

	// refund
	err = p.Refund(chargeID, amount)
	require.NoError(t, err)
}

func TestPayer_CountPrice(t *testing.T) {
	p := NewPayer(api_key)

	start := time.Now()
	end := time.Now().Add(time.Hour * 24 * 5)

	fmt.Println(start.Nanosecond(), start.Unix(), end.Nanosecond(), end.Unix())

	price := p.CountPrice(50, &models.CreateRentReq{
		RentStart: &start,
		RentEnd:   &end,
	})
	require.Equal(t, float32(250), price/100)
}
