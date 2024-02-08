package validation

import (
	"github.com/alserov/rently/proto/gen/user"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"regexp"
)

type Validator interface {
	ValidateRegisterReq(req *user.RegisterReq) error
	ValidateLoginReq(req *user.LoginReq) error
	ValidateGetInfoForRentReq(req *user.GetInfoForRentReq) error
	ValidateGetInfoReq(req *user.GetInfoReq) error
	ValidateSwitchNotificationsStatusReq(req *user.SwitchNotificationsStatusReq) error
	ValidateCheckIfAuthorizedReq(req *user.CheckIfAuthorizedReq) error
	ValidateResetPasswordReq(req *user.ResetPasswordReq) error
}

func NewValidator() Validator {
	return &validator{
		regExpEmail:    regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`),
		regExpPhone:    regexp.MustCompile(`^[1-9]\d{9}$`),
		regExpPassport: regexp.MustCompile(`^[A-Z0-9]{9}$`),
	}
}

const (
	ERR_INVALID_EMAIL           = "invalid email provided"
	ERR_INVALID_PHONE_NUMBER    = "provided invalid phone number"
	ERR_INVALID_PASSWORD        = "password is too short"
	ERR_INVALID_USERNAME        = "username can not be empty"
	ERR_INVALID_PASSPORT_NUMBER = "provided invalid passport number"
	ERR_EMPTY_UUID              = "uuid can not be empty"
	ERR_EMPTY_TOKEN             = "token can not be empty"
)

type validator struct {
	regExpEmail    *regexp.Regexp
	regExpPhone    *regexp.Regexp
	regExpPassport *regexp.Regexp
}

func (v validator) ValidateResetPasswordReq(req *user.ResetPasswordReq) error {
	if err := validateToken(req.Token); err != nil {
		return err
	}

	if err := validatePassword(req.OldPassword); err != nil {
		return err
	}

	if err := validatePassword(req.NewPassword); err != nil {
		return err
	}

	return nil
}

func (v validator) ValidateGetInfoForRentReq(req *user.GetInfoForRentReq) error {
	if err := validateToken(req.Token); err != nil {
		return err
	}

	return nil
}

func (v validator) ValidateGetInfoReq(req *user.GetInfoReq) error {
	if req.UUID == "" {
		return status.Error(codes.InvalidArgument, ERR_EMPTY_UUID)
	}

	return nil
}

func (v validator) ValidateRegisterReq(req *user.RegisterReq) error {
	if ok := v.regExpEmail.MatchString(req.Email); !ok {
		return status.Error(codes.InvalidArgument, ERR_INVALID_EMAIL)
	}

	if ok := v.regExpPhone.MatchString(req.PhoneNumber); !ok {
		return status.Error(codes.InvalidArgument, ERR_INVALID_PHONE_NUMBER)
	}

	if ok := v.regExpPassport.MatchString(req.PassportNumber); !ok {
		return status.Error(codes.InvalidArgument, ERR_INVALID_PASSPORT_NUMBER)
	}

	if err := validatePassword(req.Password); err != nil {
		return err
	}

	if len(req.Username) < 1 {
		return status.Error(codes.InvalidArgument, ERR_INVALID_USERNAME)
	}

	return nil
}

func (v validator) ValidateLoginReq(req *user.LoginReq) error {
	if ok := v.regExpEmail.MatchString(req.Email); !ok {
		return status.Error(codes.InvalidArgument, ERR_INVALID_EMAIL)
	}

	if err := validatePassword(req.Password); err != nil {
		return err
	}

	return nil
}

func (v validator) ValidateSwitchNotificationsStatusReq(req *user.SwitchNotificationsStatusReq) error {
	if req.GetUUID() == "" {
		return status.Error(codes.InvalidArgument, ERR_EMPTY_UUID)
	}

	return nil
}

func (v validator) ValidateCheckIfAuthorizedReq(req *user.CheckIfAuthorizedReq) error {
	if req.GetToken() == "" {
		return status.Error(codes.InvalidArgument, ERR_EMPTY_TOKEN)
	}

	return nil
}

func validatePassword(password string) error {
	if len(password) < 7 {
		return status.Error(codes.InvalidArgument, ERR_INVALID_PASSWORD)
	}
	return nil
}

func validateToken(token string) error {
	if token == "" {
		return status.Error(codes.InvalidArgument, ERR_EMPTY_TOKEN)
	}
	return nil
}
