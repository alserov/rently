package payment

import (
	"fmt"
	"github.com/alserov/rently/car/internal/models"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/charge"
	"github.com/stripe/stripe-go/refund"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"time"
)

type Payer interface {
	CountPrice(pricePerDay float32, service Service) float32

	Refund(chargeID string, amount float32) error
	Debit(source string, amount float32) (string, error)
}

func NewPayer(apiKey string) Payer {
	stripe.Key = apiKey
	return &payer{
		key: apiKey,
	}
}

type payer struct {
	key string
}

type Service interface {
	Period() time.Duration
}

func (p payer) CountPrice(pricePerDay float32, service Service) float32 {
	return pricePerDay * float32(service.Period().Hours()/24) * 100
}

func (p payer) Refund(chargeID string, amount float32) error {
	params := &stripe.RefundParams{
		Amount:               stripe.Int64(int64(amount)),
		Charge:               stripe.String(chargeID),
		RefundApplicationFee: stripe.Bool(false),
	}

	_, err := refund.New(params)
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	return nil
}

func (p payer) Debit(source string, amount float32) (string, error) {
	if amount < 1 {
		return "", &models.Error{
			Code: http.StatusInternalServerError,
			Msg:  fmt.Sprintf("invalid amount: %v", amount),
		}
	}

	params := &stripe.ChargeParams{
		Amount:      stripe.Int64(int64(amount)),
		Currency:    stripe.String(string(stripe.CurrencyUSD)),
		Description: stripe.String("debit card balance"),
	}

	if err := params.SetSource(source); err != nil {
		return "", &models.Error{
			Code: http.StatusInternalServerError,
			Msg:  fmt.Sprintf("failed to set source: %v", err),
		}
	}

	ch, err := charge.New(params)
	if err != nil {
		return "", &models.Error{
			Code: http.StatusInternalServerError,
			Msg:  fmt.Sprintf("failed to init charge: %v", err),
		}
	}

	return ch.ID, nil
}
