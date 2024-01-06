package service

import (
	"context"
	"fmt"
	"github.com/alserov/file-storage/internal/utils/files"
	"log/slog"
)

type Service interface {
	SetLogger(log *slog.Logger)

	GetLinks(ctx context.Context, uuid string) ([]string, error)
	GetImage(ctx context.Context, link string) ([]byte, error)
}

func NewService() Service {
	return &service{
		filer: files.NewFiler(files.RelativeImageDir),
	}
}

type service struct {
	log *slog.Logger

	filer files.Filer
}

func (s service) SetLogger(log *slog.Logger) {
	s.log = log
}

func (s service) GetLinks(ctx context.Context, uuid string) ([]string, error) {
	links, err := s.filer.GetLinks(uuid)
	if err != nil {
		return nil, fmt.Errorf("failed to get links: %w", err)
	}

	return links, nil
}

func (s service) GetImage(ctx context.Context, link string) ([]byte, error) {
	imageBytes, err := s.filer.GetImage(link)
	if err != nil {
		return nil, fmt.Errorf("failed to get image bytes: %w", err)
	}

	return imageBytes, nil
}
