package validation

import (
	"github.com/stretchr/testify/require"
	"testing"
)

type CredentialsTest struct {
	Creds string
	Valid bool
}

func TestValidator_ValidatePhoneNumber(t *testing.T) {
	tests := []CredentialsTest{
		{
			Creds: "2342423423411",
			Valid: false,
		},
		{
			Creds: "7978e433386",
			Valid: false,
		},
		{
			Creds: "",
			Valid: false,
		},
		{
			Creds: "79781433386",
			Valid: true,
		},
		{
			Creds: "+7978143338",
			Valid: false,
		},
	}

	v := NewValidator()

	for _, tc := range tests {
		err := v.ValidatePhoneNumber(tc.Creds)
		require.Equal(t, tc.Valid, err == nil, tc.Creds)
	}
}

func TestValidator_ValidateCardCredentials(t *testing.T) {
	tests := []CredentialsTest{
		{
			Creds: "1111-2222-3333-4444",
			Valid: false,
		},
		{
			Creds: "11112222333344445",
			Valid: false,
		},
		{
			Creds: "1111 2222 3333 4444",
			Valid: false,
		},
		{
			Creds: "111122223333444",
			Valid: false,
		},
		{
			Creds: "1355279860457201",
			Valid: true,
		},
	}

	v := NewValidator()

	for _, tc := range tests {
		err := v.ValidateCardCredentials(tc.Creds)
		require.Equal(t, tc.Valid, err == nil, tc.Creds)
	}
}
