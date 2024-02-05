package models

import "time"

type Error struct {
	Err string `json:"error"`
}

func (e Error) Error() string {
	return e.Err
}

type CreateRentReq struct {
	CarUUID        string `json:"carUUID" validate:"required"`
	PhoneNumber    string `json:"phoneNumber"`
	PassportNumber string `json:"passportNumber"`
	PaymentSource  string `json:"paymentSource" validate:"required"`
	RentStart      int64  `json:"rentStart"`
	RentEnd        int64  `json:"rentEnd"`
}

type CreateCarReq struct {
	Brand       string   `json:"brand" validate:"required"`
	Type        string   `json:"type" validate:"required"`
	MaxSpeed    int32    `json:"maxSpeed" validate:"required,gt=0,lt=700"`
	Seats       int32    `json:"seats" validate:"required,gt=0,lt=12"`
	Category    string   `json:"category" validate:"required"`
	PricePerDay float32  `json:"pricePerDay" validate:"required,gt=0"`
	Images      [][]byte `json:"images" validate:"required,min=1,max=7"`
	MainImage   []byte   `json:"mainImage" type:"main"  validate:"required"`
}

type UpdateCarPriceReq struct {
	UUID        string  `json:"uuid" validate:"required,min=10,max=100"`
	PricePerDay float32 `json:"pricePerDay" validate:"required,gt=0"`
}

type GetAvailableCarsReq struct {
	Start *time.Time `json:"start" validate:"required"`
	End   *time.Time `json:"end" validate:"required"`
}

type GetCarsByParamsReq struct {
	Brand       string  `json:"brand" validate:"required"`
	Type        string  `json:"type" validate:"required"`
	MaxSpeed    int32   `json:"maxSpeed" validate:"required,gt=0,lt=900"`
	Seats       int32   `json:"seats" validate:"required,gt=0,lt=12"`
	Category    string  `json:"category" validate:"required"`
	PricePerDay float32 `json:"pricePerDay" validate:"required,gt=0"`
}
