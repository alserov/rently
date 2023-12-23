package validation

import (
	"github.com/stretchr/testify/require"
	"testing"
)

type PhoneNumberTest struct {
	PhoneNumber string
	Valid       bool
}

func TestServiceValidator_ValidatePhoneNumber(t *testing.T) {
	tests := []PhoneNumberTest{
		{
			PhoneNumber: "2342423423411",
			Valid:       false,
		},
		{
			PhoneNumber: "7978e433386",
			Valid:       false,
		},
		{
			PhoneNumber: "",
			Valid:       false,
		},
		{
			PhoneNumber: "79781433386",
			Valid:       true,
		},
		{
			PhoneNumber: "+7978143338",
			Valid:       false,
		},
	}

	v := NewServiceValidator()

	for _, tc := range tests {
		err := v.ValidatePhoneNumber(tc.PhoneNumber)
		require.Equal(t, tc.Valid, err == nil, tc.PhoneNumber)
	}
}
