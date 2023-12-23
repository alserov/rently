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
}

func NewValidator() Validator {
	return &validator{
		phone: regexp.MustCompile(`\b(\d{4})\d{7}\b`),
		card:  regexp.MustCompile(`^[0-9]{13}(?:[0-9]{3})?$`),
	}
}

const (
	ERR_EMPTY                = "can not be empty"
	ERR_INVALID_PHONE_NUMBER = "invalid phone number"
	ERR_INVALID_CARD_NUMBER  = "invalid card number"
)

type validator struct {
	phone *regexp.Regexp
	card  *regexp.Regexp
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
		return status.Error(codes.InvalidArgument, "invalid rent end timestamp")
	}

	if req.GetRentStart() == nil || req.GetRentEnd().AsTime().Unix() < time.Now().Unix() {
		return status.Error(codes.InvalidArgument, "invalid rent end timestamp")
	}

	if req.GetRentEnd().AsTime().Unix() < req.GetRentStart().AsTime().Unix() {
		return status.Error(codes.InvalidArgument, "rent end can not be less than rent start")
	}

	if req.GetCardCredentials() == "" {
		return status.Error(codes.InvalidArgument, fmt.Sprintf("card number %s", ERR_EMPTY))
	}

	if req.GetCarUUID() == "" {
		return status.Error(codes.InvalidArgument, fmt.Sprintf("car uuid %s", ERR_EMPTY))
	}

	if err := v.validatePhoneNumber(req.GetPhoneNumber()); err != nil {
		return status.Error(codes.InvalidArgument, ERR_INVALID_PHONE_NUMBER)
	}

	if err := v.validateCardCredentials(req.CardCredentials); err != nil {
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
