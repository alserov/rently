package validation

type ServerValidator interface {
	ValidateCreateRentReq() error
}

func NewServerValidator() ServerValidator {
	return &serverValidator{}
}

type serverValidator struct {
}

func (v serverValidator) ValidateCreateRentReq() error {
	//TODO implement me
	panic("implement me")
}
