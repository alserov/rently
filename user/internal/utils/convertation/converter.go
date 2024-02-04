package convertation

import (
	"github.com/alserov/rently/proto/gen/user"
	"github.com/alserov/rently/user/internal/models"
)

type Converter interface {
	ToService
	ToPb
}

type ToService interface {
	RegisterReqToService(req *user.RegisterReq) models.RegisterReq
	LoginReqToService(req *user.LoginReq) models.LoginReq
}

type ToPb interface {
	RegisterResToPb(res models.RegisterRes) *user.RegisterRes
	LoginResToPb(token string) *user.LoginRes
	UserInfoResToPb(res models.UserInfoRes) *user.UserInfoRes
	InfoForRentResToPb(res models.InfoForRentRes) *user.GetInfoForRentRes
	CheckIfAuthorizedResToPb(role string) *user.CheckIfAuthorizedRes
}

func NewConverter() Converter {
	return &converter{}
}

type converter struct {
}

func (c converter) LoginResToPb(token string) *user.LoginRes {
	return &user.LoginRes{Token: token}
}

func (c converter) UserInfoResToPb(res models.UserInfoRes) *user.UserInfoRes {
	return &user.UserInfoRes{
		Username:          res.Username,
		NotificationsOn:   res.NotificationsOn,
		CurrentRentsUUIDs: res.CurrentRentsUUIDs,
	}
}

func (c converter) InfoForRentResToPb(res models.InfoForRentRes) *user.GetInfoForRentRes {
	return &user.GetInfoForRentRes{
		PhoneNumber:    res.PhoneNumber,
		PassportNumber: res.PassportNumber,
	}
}

func (c converter) CheckIfAuthorizedResToPb(role string) *user.CheckIfAuthorizedRes {
	return &user.CheckIfAuthorizedRes{
		IsAuthorized: role != "",
		Role:         role,
	}
}

func (c converter) RegisterReqToService(req *user.RegisterReq) models.RegisterReq {
	return models.RegisterReq{
		Username:       req.Username,
		Password:       req.Password,
		Email:          req.Email,
		PassportNumber: req.PassportNumber,
		PaymentSource:  req.PaymentSource,
		PhoneNumber:    req.PhoneNumber,
	}
}

func (c converter) LoginReqToService(req *user.LoginReq) models.LoginReq {
	return models.LoginReq{
		Email:    req.Email,
		Password: req.Password,
	}
}

func (c converter) RegisterResToPb(res models.RegisterRes) *user.RegisterRes {
	return &user.RegisterRes{
		UUID:  res.UUID,
		Token: res.Token,
	}
}
