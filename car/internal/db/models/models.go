package models

import "time"

type Rent struct {
	CarUUID   string
	RentPrice float32

	RentStart *time.Time
	RentEnd   *time.Time
}

type CreateRentReq struct {
	RentUUID  string  `db:"rent_uuid"`
	RentPrice float32 `db:"rent_price"`

	CarUUID string `db:"car_uuid"`

	PhoneNumber    string `db:"phone_number"`
	PassportNumber string `db:"passport_number"`
	ChargeID       string `db:"charge_id"`

	RentStart *time.Time `db:"rent_start"`
	RentEnd   *time.Time `db:"rent_end"`
}

type CancelRentInfo struct {
	RentPrice float32 `db:"rent_price"`
	ChargeID  string  `db:"chargeID"`
}

type CheckIfCarAvailable struct {
	CarUUID   string
	RentStart *time.Time
	RentEnd   *time.Time
}

// ============================================

type Car struct {
	Brand       string  `db:"brand"`
	Type        string  `db:"type"`
	MaxSpeed    int32   `db:"max_speed"`
	Seats       int32   `db:"seats"`
	Category    string  `db:"category"`
	PricePerDay float32 `db:"price_per_day"`
	UUID        string  `db:"uuid"`
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
