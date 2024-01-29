package validation

import (
	"fmt"
	"github.com/alserov/rently/proto/gen/carsharing"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"regexp"
)

type Validator interface {
	ValidateCreateRentReq(req *carsharing.CreateRentReq) error
	ValidateCancelRentReq(req *carsharing.CancelRentReq) error
	ValidateCheckRentReq(req *carsharing.CheckRentReq) error

	ValidateGetCarsByParamsReq(req *carsharing.GetCarsByParamsReq) error
	ValidateGetCarByUUID(req *carsharing.GetCarByUUIDReq) error
	ValidateGetAvailableCarsReq(req *carsharing.GetAvailableCarsReq) error
	ValidateGetCarImageReq(req *carsharing.GetImageReq) error

	ValidateCreateCarReq(req *carsharing.CreateCarReq) error
	ValidateDeleteCarReq(req *carsharing.DeleteCarReq) error
	ValidateUpdateCarPriceReq(req *carsharing.UpdateCarPriceReq) error
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
	ERR_INVALID_IMAGES_AMOUNT = "the carsharing should have at least one image"
)

type validator struct {
	phone *regexp.Regexp
	card  *regexp.Regexp
}

func (v *validator) ValidateGetCarImageReq(req *carsharing.GetImageReq) error {
	if req.GetBucket() == "" || req.GetId() == "" {
		return fmt.Errorf("image bucket or id %v", ERR_EMPTY)
	}
	return nil
}

func (v *validator) ValidateDeleteCarReq(req *carsharing.DeleteCarReq) error {
	if req.GetCarUUID() == "" {
		return status.Error(codes.InvalidArgument, fmt.Sprintf("carsharing uuid %s", ERR_EMPTY))
	}
	return nil
}

func (v *validator) ValidateUpdateCarPriceReq(req *carsharing.UpdateCarPriceReq) error {
	if req.GetCarUUID() == "" {
		return status.Error(codes.InvalidArgument, fmt.Sprintf("carsharing uuid %s", ERR_EMPTY))
	}

	if req.GetPricePerDay() <= 0 {
		return status.Error(codes.InvalidArgument, ERR_INVALID_PRICE_PER_DAY)
	}

	return nil
}

func (v *validator) ValidateCreateCarReq(req *carsharing.CreateCarReq) error {
	if req.GetBrand() == "" {
		return status.Error(codes.InvalidArgument, fmt.Sprintf("carsharing brand %s", ERR_EMPTY))
	}

	if req.GetCategory() == "" {
		return status.Error(codes.InvalidArgument, fmt.Sprintf("carsharing category %s", ERR_EMPTY))
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
		return status.Error(codes.InvalidArgument, fmt.Sprintf("carsharing type %s", ERR_EMPTY))
	}

	if len(req.GetImages()) == 0 {
		return status.Error(codes.InvalidArgument, ERR_INVALID_IMAGES_AMOUNT)
	}

	if len(req.GetMainImage()) == 0 {
		return status.Error(codes.InvalidArgument, ERR_INVALID_IMAGES_AMOUNT)
	}

	return nil
}

func (v *validator) ValidateGetCarsByParamsReq(req *carsharing.GetCarsByParamsReq) error {
	if req.GetMaxSpeed() < 0 {
		return status.Error(codes.InvalidArgument, ERR_INVALID_SPEED)
	}

	if req.GetSeats() < 0 {
		return status.Error(codes.InvalidArgument, ERR_INVALID_SEATS_AMOUNT)
	}

	if req.GetPricePerDay() < 0 {
		return status.Error(codes.InvalidArgument, ERR_INVALID_PRICE_PER_DAY)
	}

	return nil
}

func (v *validator) ValidateGetCarByUUID(req *carsharing.GetCarByUUIDReq) error {
	if req.GetUUID() == "" {
		return status.Error(codes.InvalidArgument, fmt.Sprintf("carsharing uuid %s", ERR_EMPTY))
	}

	return nil
}

func (v *validator) ValidateGetAvailableCarsReq(req *carsharing.GetAvailableCarsReq) error {
	if req.GetStart().AsTime().Before(req.GetEnd().AsTime()) {
		return status.Error(codes.InvalidArgument, "rent end can not be earlier than rent start")
	}

	return nil
}

func (v *validator) ValidateCancelRentReq(req *carsharing.CancelRentReq) error {
	if req.GetRentUUID() == "" {
		return status.Error(codes.InvalidArgument, fmt.Sprintf("rent uuid %s", ERR_EMPTY))
	}

	return nil
}

func (v *validator) ValidateCheckRentReq(req *carsharing.CheckRentReq) error {
	if req.GetRentUUID() == "" {
		return status.Error(codes.InvalidArgument, fmt.Sprintf("rent uuid %s", ERR_EMPTY))
	}

	return nil
}

func (v *validator) ValidateCreateRentReq(req *carsharing.CreateRentReq) error {
	if !req.GetRentEnd().AsTime().After(req.GetRentStart().AsTime()) {
		return status.Error(codes.InvalidArgument, "invalid rent end timestamp")
	}

	if req.GetCarUUID() == "" {
		return status.Error(codes.InvalidArgument, fmt.Sprintf("carsharing uuid %s", ERR_EMPTY))
	}

	if req.UuidIfAuthorized != "" {
		if req.GetPaymentSource() == "" {
			return status.Error(codes.InvalidArgument, fmt.Sprintf("card number %s", ERR_EMPTY))
		}

		//if err := v.validateCardCredentials(req.GetPaymentSource()); err != nil {
		//	return status.Error(codes.InvalidArgument, ERR_INVALID_CARD_NUMBER)
		//}

		if err := v.validatePhoneNumber(req.GetPhoneNumber()); err != nil {
			return status.Error(codes.InvalidArgument, ERR_INVALID_PHONE_NUMBER)
		}
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
