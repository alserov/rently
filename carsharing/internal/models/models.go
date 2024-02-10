package models

import (
	"time"
)

type Rent struct {
	CarUUID   string  `db:"car_uuid"`
	RentPrice float32 `db:"charge_amount"`
	Status    string  `db:"status"`

	RentStart time.Time `db:"rent_start"`
	RentEnd   time.Time `db:"rent_end"`
}

type RentStartData struct {
	CarUUID   string    `db:"car_uuid"`
	UserUUID  string    `db:"user_uuid"`
	RentStart time.Time `db:"rent_start"`
	RentEnd   time.Time `db:"rent_end"`
}

type CreateRentReq struct {
	RentUUID string
	ChargeID string
	UserUUID string

	Token string

	CarUUID        string
	PhoneNumber    string
	PassportNumber string
	PaymentSource  string

	RentStart time.Time
	RentEnd   time.Time
	Days      int
}

type CancelRentInfo struct {
	ChargeID  string  `db:"uuid"`
	RentPrice float32 `db:"rent_price"`
}

type Charge struct {
	RentUUID     string  `db:"rent_uuid"`
	ChargeUUID   string  `db:"charge_uuid"`
	ChargeAmount float32 `db:"charge_amount"`
}

func (c *CreateRentReq) Period() time.Duration {
	period := c.RentEnd.Sub(c.RentStart)
	return period
}

type CreateRentRes struct {
	RentUUID string
}

// ============================================

type CarMainInfo struct {
	UUID        string  `db:"uuid"`
	Brand       string  `db:"brand"`
	Type        string  `db:"type"`
	Category    string  `db:"category"`
	PricePerDay float32 `db:"price_per_day"`
	Image       string  `db:"image"`
}

type Car struct {
	UUID        string `db:"uuid"`
	MainImage   string
	Images      []string
	Brand       string  `db:"brand"`
	Type        string  `db:"type"`
	MaxSpeed    int32   `db:"max_speed"`
	Seats       int32   `db:"seats"`
	Category    string  `db:"category"`
	PricePerDay float32 `db:"price_per_day"`
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

type UpdateCarPriceReq struct {
	CarUUID string
	Price   float32
}

type UserInfo struct {
	PhoneNumber    string
	PassportNumber string
	UUID           string
}
