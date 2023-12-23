package models

import "time"

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

type Rent struct {
	CarUUID   string
	RentPrice float32

	RentStart *time.Time
	RentEnd   *time.Time
}

type CancelRentInfo struct {
	RentPrice       float32 `db:"rent_price"`
	CardCredentials string  `db:"card_credentials"`
}
