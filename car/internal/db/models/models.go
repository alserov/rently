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

	PhoneNumber     string `db:"phone_number"`
	PassportNumber  string `db:"passport_number"`
	CardCredentials string `db:"card_credentials"`

	RentStart *time.Time `db:"rent_start"`
	RentEnd   *time.Time `db:"rent_end"`
}

type CancelRentInfo struct {
	RentPrice       float32 `db:"rent_price"`
	CardCredentials string  `db:"card_credentials"`
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
