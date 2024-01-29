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
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	Register(ctx context.Context, req models.RegisterReq) (models.RegisterRes, error)
	Login(ctx context.Context, req models.LoginReq) (string, error)
	GetInfo(ctx context.Context, uuid string) (models.UserInfoRes, error)
	GetRentInfo(ctx context.Context, uuid string) (models.InfoForRentRes, error)
	SwitchNotificationsStatus(ctx context.Context, uuid string) error
	CheckIfAuthorized(ctx context.Context, token string) (string, error)
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

type service struct {
	log log.Logger

	notifier notifications.Notifier

	repo db.Repository
}

func (s *service) CheckIfAuthorized(ctx context.Context, token string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (s *service) GetInfo(ctx context.Context, uuid string) (models.UserInfoRes, error) {
	info, err := s.repo.GetUserInfo(ctx, uuid)
	if err != nil {
		return models.UserInfoRes{}, nil
	}

	return info, nil
}

func (s *service) GetRentInfo(ctx context.Context, uuid string) (models.InfoForRentRes, error) {
	info, err := s.repo.GetInfoForRent(ctx, uuid)
	if err != nil {
		return models.InfoForRentRes{}, err
	}

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

	if err = bcrypt.CompareHashAndPassword([]byte(userData.Password), []byte(req.Password)); err != nil {
		return "", fmt.Errorf("invalid password")
	}

	token, err := newToken(userData.UUID)
	if err != nil {
		return "", fmt.Errorf("failed to generate new token: %w", err)
	}

	if err = s.notifier.Login(userData.Email); err != nil {
		return "", fmt.Errorf("failed to send login notification: %w", err)
	}

	return token, nil
}

func (s *service) Register(ctx context.Context, req models.RegisterReq) (models.RegisterRes, error) {
	req.UUID = uuid.New().String()

	hashedPassword, err := hash(req.Password)
	if err != nil {
		return models.RegisterRes{}, fmt.Errorf("failed to hash password: %w", err)
	}

	req.Password = hashedPassword
	req.PassportNumber = encrypt(req.PassportNumber)
	req.PaymentSource = encrypt(req.PaymentSource)

	if err = s.repo.Register(ctx, req); err != nil {
		return models.RegisterRes{}, err
	}

	if err = s.notifier.Registration(req.Email); err != nil {
		return models.RegisterRes{}, fmt.Errorf("failed to send notification: %w", err)
	}

	token, err := newToken(req.UUID)
	if err != nil {
		return models.RegisterRes{}, fmt.Errorf("failed to generate new token: %w", err)
	}

	return models.RegisterRes{
		UUID:  req.UUID,
		Token: token,
	}, nil
}
