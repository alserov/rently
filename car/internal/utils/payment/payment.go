package payment

import "time"

type Payer interface {
	CountPrice(service Service) float32

	TopUp(cardNumber string, amount float32) error
	Debit(cardNumber string, amount float32) error
}

func NewPayer(apiKey string) Payer {
	return &payer{
		key: apiKey,
	}
}

type payer struct {
	key string
}

type Service interface {
	Price() float32
	Period() time.Duration
}

func (p payer) CountPrice(service Service) float32 {
	return service.Price() * float32(service.Period())
}

func (p payer) TopUp(cardNumber string, amount float32) error {
	panic("")
}

func (p payer) Debit(cardNumber string, amount float32) error {
	//TODO implement me
	panic("implement me")
}
