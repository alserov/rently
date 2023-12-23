package validation

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"regexp"
)

type ServiceValidator interface {
	ValidatePhoneNumber(phoneNumber string) error
	ValidateCardCredentials(cardCredentials string) error
	ValidatePassportNumber(passportNumber string) error
}

func NewServiceValidator() ServiceValidator {
	return &serviceValidator{
		phone: regexp.MustCompile(`\b(\d{4})\d{7}\b`),
	}
}

type serviceValidator struct {
	phone *regexp.Regexp
}

const (
	ERR_INVALID_PHONE_NUMBER = "invalid phone number"
)

func (s serviceValidator) ValidatePhoneNumber(phoneNumber string) error {
	correct := s.phone.MatchString(phoneNumber)
	if !correct {
		return status.Error(codes.InvalidArgument, ERR_INVALID_PHONE_NUMBER)
	}

	return nil
}

func (s serviceValidator) ValidateCardCredentials(cardCredentials string) error {
	//TODO implement me
	panic("implement me")
}

func (s serviceValidator) ValidatePassportNumber(passportNumber string) error {
	//TODO implement me
	panic("implement me")
}
