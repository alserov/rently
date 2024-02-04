package convertation

import (
	"fmt"
	"github.com/alserov/rently/proto/gen/carsharing"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
	"testing"
	"time"
)

func TestServerConverter_CreateRentReqToService(t *testing.T) {
	c := NewServerConverter()

	now := time.Now()

	rentStartPb := timestamppb.Timestamp{
		Seconds: now.Unix(),
		Nanos:   int32(now.Nanosecond()),
	}
	rentEndPb := timestamppb.Timestamp{
		Seconds: now.Add(time.Hour * 24).Unix(),
		Nanos:   int32(now.Add(time.Hour * 24).Nanosecond()),
	}

	fmt.Println(rentStartPb.AsTime(), rentEndPb.AsTime())

	req := carsharing.CreateRentReq{
		CarUUID:        "uuid",
		PhoneNumber:    "42423423424",
		PassportNumber: "2342424",
		PaymentSource:  "1234131313",
		RentStart:      &rentStartPb,
		RentEnd:        &rentEndPb,
	}

	converted := c.CreateRentReqToService(&req)
	require.Equal(t, req.CarUUID, converted.CarUUID, "carsharing uuid")
	require.Equal(t, req.PhoneNumber, converted.PhoneNumber, "phone number")
	require.Equal(t, req.PaymentSource, converted.PaymentSource, "payment source")
	require.Equal(t, req.PassportNumber, converted.PassportNumber, "passport number")
	require.Equal(t, req.RentStart.AsTime(), converted.RentStart, "rent start")
	require.Equal(t, req.RentEnd.AsTime(), converted.RentEnd, "rent end")
}
