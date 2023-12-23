package validation

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"regexp"
)

type Validator interface {
	ValidatePhoneNumber(phoneNumber string) error
	ValidateCardCredentials(cardCredentials string) error
	ValidatePassportNumber(passportNumber string) error

	ValidateCreateRentReq() error
}

func NewValidator() Validator {
	return &validator{
		phone: regexp.MustCompile(`\b(\d{4})\d{7}\b`),
		card:  regexp.MustCompile(`^[0-9]{13}(?:[0-9]{3})?$`),
	}
}

type validator struct {
	phone *regexp.Regexp
	card  *regexp.Regexp
}

func (s *validator) ValidateCreateRentReq() error {
	//TODO implement me
	panic("implement me")
}

const (
	ERR_INVALID_PHONE_NUMBER = "invalid phone number"
	ERR_INVALID_CARD_NUMBER  = "invalid card number"
)

func (s *validator) ValidatePhoneNumber(phoneNumber string) error {
	valid := s.phone.MatchString(phoneNumber)
	if !valid {
		return status.Error(codes.InvalidArgument, ERR_INVALID_PHONE_NUMBER)
	}

	return nil
}

func (s *validator) ValidateCardCredentials(cardCredentials string) error {
	valid := s.card.MatchString(cardCredentials)
	if !valid {
		return status.Error(codes.InvalidArgument, ERR_INVALID_CARD_NUMBER)
	}

	return nil
}

func (s *validator) ValidatePassportNumber(passportNumber string) error {
	//TODO implement me
	panic("implement me")
}
