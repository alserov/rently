package models

import "time"

type Rent struct {
	CarUUID   string
	RentPrice float32

	RentStart *time.Time
	RentEnd   *time.Time
}

type CreateRentReq struct {
	RentUUID string

	CarPricePerDay float32
	CarUUID        string
	PhoneNumber    string
	PassportNumber string
	PaymentSource  string

	RentStart *time.Time
	RentEnd   *time.Time
}

func (c *CreateRentReq) Period() time.Duration {
	period := c.RentEnd.Sub(*c.RentStart)
	return period
}

type CreateRentRes struct {
	RentUUID string
	ChargeID string
}

type CancelRentReq struct {
	RentUUID string
}

type CheckRentReq struct {
	RentUUID string
}

// ============================================

type Car struct {
	Brand       string
	Type        string
	MaxSpeed    int32
	Seats       int32
	Category    string
	PricePerDay float32
	UUID        string
}

type Period struct {
	Start *time.Time
	End   *time.Time
}

type CarParams struct {
	Brand       string
	Type        string
	MaxSpeed    int32
	Seats       int32
	Category    string
	PricePerDay float32
}
