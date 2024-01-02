package payment

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestPayer(t *testing.T) {
	p := NewPayer("sk_test_51OU56CDOnc0MdcTNBwddO2cn8NrEebjfuAGjBjj9xSyKmiUO4ajJ1vZ0yBoOsAMq0HjHqCmis2niwoj2EZYCDLOA00lcCUlWxh")

	// invalid price
	chargeID, err := p.Debit(-10000)
	require.Error(t, err)
	require.Empty(t, chargeID)

	const amount = 100_00
	// valid price
	chargeID, err = p.Debit(amount)
	require.NoError(t, err)
	require.NotEmpty(t, chargeID)

	// refund
	err = p.Refund(chargeID, amount)
	require.NoError(t, err)
}
