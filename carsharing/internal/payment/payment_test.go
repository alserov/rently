package payment

import (
	"github.com/stretchr/testify/require"
	"testing"
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
