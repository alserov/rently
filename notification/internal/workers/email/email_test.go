package email

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestEmailer_Send(t *testing.T) {
	e := NewEmailer(Params{
		From:     "asroka4ev@gmail.com",
		Password: "bwch rbes cemw ayyc",
		SmtpHost: "smtp.gmail.com",
		SmtpPort: 587,
	})
	err := e.Registration("myweirdmailik@gmail.com")
	require.NoError(t, err)
}
