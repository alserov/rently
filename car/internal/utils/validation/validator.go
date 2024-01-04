package validation

import (
	"fmt"
	"github.com/alserov/rently/proto/gen/car"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"regexp"
	"time"
)

type Validator interface {
	ValidateCreateRentReq(req *car.CreateRentReq) error
	ValidateCancelRentReq(req *car.CancelRentReq) error
	ValidateCheckRentReq(req *car.CheckRentReq) error

	ValidateGetCarsByParamsReq(req *car.GetCarsByParamsReq) error
	ValidateGetCarByUUID(req *car.GetCarByUUIDReq) error
	ValidateGetAvailableCarsReq(req *car.GetAvailableCarsReq) error

	ValidateCreateCarReq(req *car.CreateCarReq) error
	ValidateDeleteCarReq(req *car.DeleteCarReq) error
	ValidateUpdateCarPriceReq(req *car.UpdateCarPriceReq) error
}

func NewValidator() Validator {
	return &validator{
		phone: regexp.MustCompile(`\b(\d{4})\d{7}\b`),
		card:  regexp.MustCompile(`^[0-9]{13}(?:[0-9]{3})?$`),
	}
}

const (
	ERR_EMPTY                 = "can not be empty"
	ERR_INVALID_PHONE_NUMBER  = "invalid phone number"
	ERR_INVALID_CARD_NUMBER   = "invalid card number"
	ERR_INVALID_SPEED         = "invalid speed"
	ERR_INVALID_SEATS_AMOUNT  = "invalid seats amount"
	ERR_INVALID_PRICE_PER_DAY = "price can not be less or equal to 0.txt"
	ERR_INVALID_IMAGES_AMOUNT = "the car should have at least one image"
)

type validator struct {
	phone *regexp.Regexp
	card  *regexp.Regexp
}

func (v *validator) ValidateDeleteCarReq(req *car.DeleteCarReq) error {
	if req.GetCarUUID() == "" {
		return status.Error(codes.InvalidArgument, fmt.Sprintf("car uuid %s", ERR_EMPTY))
	}
	return nil
}

func (v *validator) ValidateUpdateCarPriceReq(req *car.UpdateCarPriceReq) error {
	if req.GetCarUUID() == "" {
		return status.Error(codes.InvalidArgument, fmt.Sprintf("car uuid %s", ERR_EMPTY))
	}

	if req.GetPricePerDay() <= 0 {
		return status.Error(codes.InvalidArgument, ERR_INVALID_PRICE_PER_DAY)
	}

	return nil
}

func (v *validator) ValidateCreateCarReq(req *car.CreateCarReq) error {
	if req.GetBrand() == "" {
		return status.Error(codes.InvalidArgument, fmt.Sprintf("car brand %s", ERR_EMPTY))
	}

	if req.GetCategory() == "" {
		return status.Error(codes.InvalidArgument, fmt.Sprintf("car category %s", ERR_EMPTY))
	}

	if req.GetSeats() < 1 {
		return status.Error(codes.InvalidArgument, ERR_INVALID_SEATS_AMOUNT)
	}

	if req.GetMaxSpeed() < 0 {
		return status.Error(codes.InvalidArgument, ERR_INVALID_SPEED)
	}

	if req.GetPricePerDay() <= 0 {
		return status.Error(codes.InvalidArgument, ERR_INVALID_PRICE_PER_DAY)
	}

	if req.GetType() == "" {
		return status.Error(codes.InvalidArgument, fmt.Sprintf("car type %s", ERR_EMPTY))
	}

	if len(req.GetImages()) == 0 {
		return status.Error(codes.InvalidArgument, ERR_INVALID_IMAGES_AMOUNT)
	}

	return nil
}

func (v *validator) ValidateGetCarsByParamsReq(req *car.GetCarsByParamsReq) error {
	if req.GetMaxSpeed() < 0 {
		return status.Error(codes.InvalidArgument, ERR_INVALID_SPEED)
	}

	if req.GetSeats() < 2 {
		return status.Error(codes.InvalidArgument, ERR_INVALID_SEATS_AMOUNT)
	}

	if req.GetPricePerDay() <= 0 {
		return status.Error(codes.InvalidArgument, ERR_INVALID_PRICE_PER_DAY)
	}

	return nil
}

func (v *validator) ValidateGetCarByUUID(req *car.GetCarByUUIDReq) error {
	if req.GetUUID() == "" {
		return status.Error(codes.InvalidArgument, fmt.Sprintf("car uuid %s", ERR_EMPTY))
	}

	return nil
}

func (v *validator) ValidateGetAvailableCarsReq(req *car.GetAvailableCarsReq) error {
	if req.GetStart().AsTime().Before(req.GetEnd().AsTime()) {
		return status.Error(codes.InvalidArgument, "rent end can not be earlier than rent start")
	}

	return nil
}

func (v *validator) ValidateCancelRentReq(req *car.CancelRentReq) error {
	if req.GetRentUUID() == "" {
		return status.Error(codes.InvalidArgument, fmt.Sprintf("rent uuid %s", ERR_EMPTY))
	}

	return nil
}

func (v *validator) ValidateCheckRentReq(req *car.CheckRentReq) error {
	if req.GetRentUUID() == "" {
		return status.Error(codes.InvalidArgument, fmt.Sprintf("rent uuid %s", ERR_EMPTY))
	}

	return nil
}

func (v *validator) ValidateCreateRentReq(req *car.CreateRentReq) error {
	if req.GetRentEnd() == nil || req.GetRentEnd().AsTime().Unix() < time.Now().Unix() {
		return status.Error(codes.InvalidArgument, "invalid car end timestamp")
	}

	if req.GetRentStart() == nil || req.GetRentEnd().AsTime().Unix() < time.Now().Unix() {
		return status.Error(codes.InvalidArgument, "invalid car end timestamp")
	}

	if req.GetRentEnd().AsTime().Unix() <= req.GetRentStart().AsTime().Unix() {
		return status.Error(codes.InvalidArgument, "rent end can not be earlier than rent start")
	}

	if req.GetPaymentSource() == "" {
		return status.Error(codes.InvalidArgument, fmt.Sprintf("card number %s", ERR_EMPTY))
	}

	if req.GetCarUUID() == "" {
		return status.Error(codes.InvalidArgument, fmt.Sprintf("car uuid %s", ERR_EMPTY))
	}

	if err := v.validatePhoneNumber(req.GetPhoneNumber()); err != nil {
		return status.Error(codes.InvalidArgument, ERR_INVALID_PHONE_NUMBER)
	}

	if err := v.validateCardCredentials(req.GetPaymentSource()); err != nil {
		return status.Error(codes.InvalidArgument, ERR_INVALID_CARD_NUMBER)
	}

	return nil
}

func (v *validator) validatePhoneNumber(phoneNumber string) error {
	valid := v.phone.MatchString(phoneNumber)
	if !valid {
		return status.Error(codes.InvalidArgument, ERR_INVALID_PHONE_NUMBER)
	}

	return nil
}

func (v *validator) validateCardCredentials(cardCredentials string) error {
	valid := v.card.MatchString(cardCredentials)
	if !valid {
		return status.Error(codes.InvalidArgument, ERR_INVALID_CARD_NUMBER)
	}

	return nil
}

func (v *validator) validatePassportNumber(passportNumber string) error {
	//TODO implement me
	panic("implement me")
}
