package validation

import "github.com/alserov/rently/proto/gen/user"

type Validator interface {
	ValidateRegisterReq(req *user.RegisterReq) error
	ValidateLoginReq(req *user.LoginReq) error
	ValidateGetInfoReq(req *user.GetInfoReq) error
	ValidateSwitchNotificationsStatusReq(req *user.SwitchNotificationsStatusReq) error
	ValidateCheckIfAuthorizedReq(req *user.CheckIfAuthorizedReq) error
}

func NewValidator() Validator {
	return &validator{}
}

type validator struct {
}

func (v validator) ValidateRegisterReq(req *user.RegisterReq) error {
	//TODO implement me
	panic("implement me")
}

func (v validator) ValidateLoginReq(req *user.LoginReq) error {
	//TODO implement me
	panic("implement me")
}

func (v validator) ValidateGetInfoReq(req *user.GetInfoReq) error {
	//TODO implement me
	panic("implement me")
}

func (v validator) ValidateSwitchNotificationsStatusReq(req *user.SwitchNotificationsStatusReq) error {
	//TODO implement me
	panic("implement me")
}

func (v validator) ValidateCheckIfAuthorizedReq(req *user.CheckIfAuthorizedReq) error {
	//TODO implement me
	panic("implement me")
}
