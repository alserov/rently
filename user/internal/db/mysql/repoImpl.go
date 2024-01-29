package mysql

import (
	"context"
	"github.com/alserov/rently/user/internal/db"
	"github.com/alserov/rently/user/internal/models"
	"github.com/jmoiron/sqlx"
)

func NewRepository(db *sqlx.DB) db.Repository {
	return &repository{}
}

type repository struct {
}

func (r repository) Register(ctx context.Context, req models.RegisterReq) error {
	//TODO implement me
	panic("implement me")
}
