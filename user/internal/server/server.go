package server

import (
	"context"
	"github.com/alserov/rently/proto/gen/user"
	"github.com/alserov/rently/user/internal/log"
	"github.com/alserov/rently/user/internal/service"
	"github.com/alserov/rently/user/internal/utils/convertation"
	"github.com/alserov/rently/user/internal/utils/validation"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Params struct {
	Service service.Service
}

func RegisterGRPCServer(s *grpc.Server, params Params) {
	user.RegisterUserServer(s, newServer(params))
}

func newServer(s Params) user.UserServer {
	return &server{
		log:     log.GetLogger(),
		service: s.Service,
		valid:   validation.NewValidator(),
		convert: convertation.NewConverter(),
	}
}

type server struct {
	user.UnimplementedUserServer

	log log.Logger

	service service.Service

	valid validation.Validator

	convert convertation.Converter
}

func (s *server) CheckIfAuthorized(ctx context.Context, req *user.CheckIfAuthorizedReq) (*user.CheckIfAuthorizedRes, error) {
	if err := s.valid.ValidateCheckIfAuthorizedReq(req); err != nil {
		return nil, err
	}

	role, err := s.service.CheckIfAuthorized(ctx, req.Token)
	if err != nil {
		return nil, s.handleError(err)
	}

	return s.convert.CheckIfAuthorizedResToPb(role), nil
}

func (s *server) Register(ctx context.Context, req *user.RegisterReq) (*user.RegisterRes, error) {
	if err := s.valid.ValidateRegisterReq(req); err != nil {
		return nil, err
	}

	res, err := s.service.Register(ctx, s.convert.RegisterReqToService(req))
	if err != nil {
		return nil, s.handleError(err)
	}

	return s.convert.RegisterResToPb(res), nil
}

func (s *server) Login(ctx context.Context, req *user.LoginReq) (*user.LoginRes, error) {
	if err := s.valid.ValidateLoginReq(req); err != nil {
		return nil, err
	}

	token, err := s.service.Login(ctx, s.convert.LoginReqToService(req))
	if err != nil {
		return nil, s.handleError(err)
	}

	return s.convert.LoginResToPb(token), nil
}

func (s *server) GetInfo(ctx context.Context, req *user.GetInfoReq) (*user.UserInfoRes, error) {
	if err := s.valid.ValidateGetInfoReq(req); err != nil {
		return nil, err
	}

	info, err := s.service.GetInfo(ctx, req.UUID)
	if err != nil {
		return nil, s.handleError(err)
	}

	return s.convert.UserInfoResToPb(info), nil
}

func (s *server) GetInfoForRent(ctx context.Context, req *user.GetInfoReq) (*user.GetInfoForRentRes, error) {
	if err := s.valid.ValidateGetInfoReq(req); err != nil {
		return nil, err
	}

	info, err := s.service.GetRentInfo(ctx, req.UUID)
	if err != nil {
		return nil, s.handleError(err)
	}

	return s.convert.InfoForRentResToPb(info), nil
}

func (s *server) SwitchStatusNotifications(ctx context.Context, req *user.SwitchNotificationsStatusReq) (*emptypb.Empty, error) {
	if err := s.valid.ValidateSwitchNotificationsStatusReq(req); err != nil {
		return nil, err
	}

	if err := s.service.SwitchNotificationsStatus(ctx, req.UUID); err != nil {
		return nil, s.handleError(err)
	}

	return &emptypb.Empty{}, nil
}
