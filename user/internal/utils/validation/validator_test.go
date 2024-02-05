package validation

import (
	"github.com/alserov/rently/proto/gen/user"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/status"
	"testing"
)

//func TestTimer(t *testing.T) {
//	n := time.Now().Add(time.Hour * 24 * 3)
//	fmt.Println(n.Unix())
//	fmt.Println(n.Add(time.Hour * 24 * 7).Unix())
//}

func TestValidator_ValidateRegisterReq(t *testing.T) {
	tests := []struct {
		name     string
		in       *user.RegisterReq
		expError string
	}{
		{
			name: "invalid email",
			in: &user.RegisterReq{
				Username:       "user123",
				Password:       "qwerty123",
				Email:          "myemail@gmail",
				PassportNumber: "AB1234567",
				PaymentSource:  "1122 3344 5566 7788",
				PhoneNumber:    "3801231223",
			},
			expError: ERR_INVALID_EMAIL,
		},
		{
			name: "invalid password",
			in: &user.RegisterReq{
				Username:       "user123",
				Password:       "qwerty",
				Email:          "myemail@gmail.com",
				PassportNumber: "AB1234567",
				PaymentSource:  "1122 3344 5566 7788",
				PhoneNumber:    "3801231223",
			},
			expError: ERR_INVALID_PASSWORD,
		},
		{
			name: "invalid passport number",
			in: &user.RegisterReq{
				Username:       "user123",
				Password:       "qwerty123",
				Email:          "myemail@gmail.com",
				PassportNumber: "AB123",
				PaymentSource:  "1122 3344 5566 7788",
				PhoneNumber:    "3801231223",
			},
			expError: ERR_INVALID_PASSPORT_NUMBER,
		},
		{
			name: "valid phone number",
			in: &user.RegisterReq{
				Username:       "user123",
				Password:       "qwerty123",
				Email:          "myemail@gmail.com",
				PassportNumber: "DE9876543",
				PaymentSource:  "1122 3344 5566 7788",
				PhoneNumber:    "380123",
			},
			expError: ERR_INVALID_PHONE_NUMBER,
		},
		{
			name: "valid",
			in: &user.RegisterReq{
				Username:       "user123",
				Password:       "qwerty123",
				Email:          "myemail@gmail.com",
				PassportNumber: "DE9876543",
				PaymentSource:  "1122 3344 5566 7788",
				PhoneNumber:    "3801232323",
			},
			expError: "",
		},
	}

	v := NewValidator()
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := v.ValidateRegisterReq(tc.in)

			st, _ := status.FromError(err)

			require.Equal(t, tc.expError, st.Message())
		})
	}
}
