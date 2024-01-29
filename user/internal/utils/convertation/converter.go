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
	//TODO implement me
	panic("implement me")
}

func (c converter) UserInfoResToPb(res models.UserInfoRes) *user.UserInfoRes {
	//TODO implement me
	panic("implement me")
}

func (c converter) InfoForRentResToPb(res models.InfoForRentRes) *user.GetInfoForRentRes {
	//TODO implement me
	panic("implement me")
}

func (c converter) CheckIfAuthorizedResToPb(role string) *user.CheckIfAuthorizedRes {
	//TODO implement me
	panic("implement me")
}

func (c converter) RegisterReqToService(req *user.RegisterReq) models.RegisterReq {
	//TODO implement me
	panic("implement me")
}

func (c converter) LoginReqToService(req *user.LoginReq) models.LoginReq {
	//TODO implement me
	panic("implement me")
}

func (c converter) RegisterResToPb(res models.RegisterRes) *user.RegisterRes {
	//TODO implement me
	panic("implement me")
}
