package service

import (
	"context"
	"fmt"
	"github.com/alserov/rently/user/internal/db"
	"github.com/alserov/rently/user/internal/log"
	"github.com/alserov/rently/user/internal/metrics"
	"github.com/alserov/rently/user/internal/models"
	"github.com/alserov/rently/user/internal/notifications"
	"github.com/google/uuid"
	"log/slog"
	"net/http"
)

type Service interface {
	Register(ctx context.Context, req models.RegisterReq) (models.RegisterRes, error)
	Login(ctx context.Context, req models.LoginReq) (string, error)
	GetInfo(ctx context.Context, uuid string) (models.UserInfoRes, error)
	GetRentInfo(ctx context.Context, token string) (models.InfoForRentRes, error)
	SwitchNotificationsStatus(ctx context.Context, uuid string) error
	CheckIfAuthorized(ctx context.Context, token string) (string, error)
	ResetPassword(ctx context.Context, req models.ResetPasswordReq) error
}

type Params struct {
	Metrics  metrics.Metrics
	Notifier notifications.Notifier
	Repo     db.Repository
}

func NewService(p Params) Service {
	return &service{
		log:      log.GetLogger(),
		notifier: p.Notifier,
		repo:     p.Repo,
	}
}

const (
	ROLE_USER  = "user"
	ROLE_ADMIN = "admin"
)

type service struct {
	log log.Logger

	notifier notifications.Notifier

	repo db.Repository
}

func (s *service) ResetPassword(ctx context.Context, req models.ResetPasswordReq) error {
	uuid, _, err := parseTokenClaims(req.Token)
	if err != nil {
		return fmt.Errorf("failed to parse token: %w", err)
	}

	passwd, err := s.repo.GetPassword(ctx, uuid)
	if err != nil {
		return err
	}

	if err = compareHashAndPassword(passwd, req.OldPassword); err != nil {
		return err
	}

	hashedPassword, err := hash(req.NewPassword)
	if err != nil {
		return err
	}

	if err = s.repo.ResetPassword(ctx, uuid, hashedPassword); err != nil {
		return err
	}

	return nil
}

func (s *service) CheckIfAuthorized(ctx context.Context, token string) (string, error) {
	uuid, role, err := parseTokenClaims(token)
	if err != nil {
		return "", fmt.Errorf("failed to parse token: %w", err)
	}

	if err = s.repo.CheckIfAuthorized(ctx, uuid, role); err != nil {
		return "", err
	}

	return role, nil
}

func (s *service) GetInfo(ctx context.Context, uuid string) (models.UserInfoRes, error) {
	info, err := s.repo.GetUserInfo(ctx, uuid)
	if err != nil {
		return models.UserInfoRes{}, nil
	}

	return info, nil
}

func (s *service) GetRentInfo(ctx context.Context, token string) (models.InfoForRentRes, error) {
	uuid, _, err := parseTokenClaims(token)
	if err != nil {
		return models.InfoForRentRes{}, err
	}

	s.log.Debug("parsed token", slog.String("uuid", uuid))

	info, err := s.repo.GetInfoForRent(ctx, uuid)
	if err != nil {
		return models.InfoForRentRes{}, err
	}

	info.UUID = uuid

	s.log.Debug("got info for rent", slog.Any("info", info))

	passportNumber, err := decrypt(info.PassportNumber)
	if err != nil {
		return models.InfoForRentRes{}, fmt.Errorf("failed to decrypt passport number: %w", err)
	}

	phoneNumber, err := decrypt(info.PhoneNumber)
	if err != nil {
		return models.InfoForRentRes{}, fmt.Errorf("failed to decrypt phone number: %w", err)
	}

	return models.InfoForRentRes{
		PhoneNumber:    phoneNumber,
		PassportNumber: passportNumber,
		UUID:           info.UUID,
	}, nil
}

func (s *service) SwitchNotificationsStatus(ctx context.Context, uuid string) error {
	if err := s.repo.SwitchNotificationsStatus(ctx, uuid); err != nil {
		return err
	}

	return nil
}

func (s *service) Login(ctx context.Context, req models.LoginReq) (string, error) {
	userData, err := s.repo.Login(ctx, req.Email)
	if err != nil {
		return "", err
	}

	if err = compareHashAndPassword(userData.Password, req.Password); err != nil {
		return "", err
	}

	token, err := newToken(userData.UUID, userData.Role)
	if err != nil {
		return "", fmt.Errorf("failed to generate new token: %w", err)
	}

	if err = s.notifier.Login(ctx, userData.Email); err != nil {
		return "", fmt.Errorf("failed to send login notification: %w", err)
	}

	return token, nil
}

func (s *service) Register(ctx context.Context, req models.RegisterReq) (models.RegisterRes, error) {
	req.UUID = uuid.New().String()

	_, err := s.repo.Login(ctx, req.Email)
	if err == nil {
		return models.RegisterRes{}, &models.Error{
			Msg:    fmt.Sprintf("user with email: %s already exists", req.Email),
			Status: http.StatusBadRequest,
		}
	}

	hashedPassword, err := hash(req.Password)
	if err != nil {
		return models.RegisterRes{}, fmt.Errorf("failed to hash password: %w", err)
	}

	req.Password = hashedPassword
	req.PassportNumber = encrypt(req.PassportNumber)
	req.PaymentSource = encrypt(req.PaymentSource)
	req.PhoneNumber = encrypt(req.PhoneNumber)

	if err = s.repo.Register(ctx, req); err != nil {
		return models.RegisterRes{}, err
	}

	if err = s.notifier.Registration(ctx, req.Email); err != nil {
		return models.RegisterRes{}, fmt.Errorf("failed to send notification: %w", err)
	}

	token, err := newToken(req.UUID, ROLE_USER)
	if err != nil {
		return models.RegisterRes{}, fmt.Errorf("failed to generate new token: %w", err)
	}

	return models.RegisterRes{
		UUID:  req.UUID,
		Token: token,
	}, nil
}
