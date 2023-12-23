package models

import "time"

type CreateRentReq struct {
	RentUUID string

	CarUUID         string
	PhoneNumber     string
	PassportNumber  string
	CardCredentials string

	RentStart *time.Time
	RentEnd   *time.Time
}

type CancelRentReq struct {
	RentUUID string
}

type CheckRentReq struct {
	RentUUID string
}

type Rent struct {
	CarUUID   string
	RentPrice float32

	RentStart *time.Time
	RentEnd   *time.Time
}
