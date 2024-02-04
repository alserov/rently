package db

import (
	"context"
	"github.com/alserov/rently/user/internal/models"
)

type Repository interface {
	Register(ctx context.Context, req models.RegisterReq) error
	Login(ctx context.Context, email string) (LoginInfo, error)
	GetUserInfo(ctx context.Context, uuid string) (models.UserInfoRes, error)
	GetInfoForRent(ctx context.Context, uuid string) (models.InfoForRentRes, error)
	SwitchNotificationsStatus(ctx context.Context, uuid string) error
	GetUserByUUID(ctx context.Context, uuids string) (EmailNotificationsInfo, error)
	CheckIfAuthorized(ctx context.Context, uuid string, role string) error
}

type LoginInfo struct {
	UUID     string
	Email    string
	Password string
}

type EmailNotificationsInfo struct {
	Email    string
	Username string
}
