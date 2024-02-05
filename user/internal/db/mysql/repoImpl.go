package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/alserov/rently/user/internal/db"
	"github.com/alserov/rently/user/internal/models"
	"github.com/jmoiron/sqlx"
	"net/http"
)

func NewRepository(db *sqlx.DB) db.Repository {
	return &repository{
		db: db,
	}
}

type repository struct {
	db *sqlx.DB
}

func (r repository) CheckIfAuthorized(ctx context.Context, uuid string, role string) error {
	query := `SELECT count(*) FROM users WHERE uuid = ?`

	err := r.db.QueryRowx(query, uuid).Err()
	if errors.Is(err, sql.ErrNoRows) {
		return &models.Error{
			Msg:    fmt.Sprintf("%s not found", role),
			Status: http.StatusNotFound,
		}
	}
	if err != nil {
		return &models.Error{
			Msg:    fmt.Sprintf("failed to find %s by uuid: %v", role, err),
			Status: http.StatusInternalServerError,
		}
	}

	return nil
}

func (r repository) Login(_ context.Context, email string) (db.LoginInfo, error) {
	query := `SELECT uuid,email,password,role FROM users WHERE email = ?  LIMIT 1`

	var info db.LoginInfo
	err := r.db.QueryRowx(query, email).StructScan(&info)
	if errors.Is(sql.ErrNoRows, err) {
		return db.LoginInfo{}, &models.Error{
			Msg:    fmt.Sprintf("user with email: %s not found", email),
			Status: http.StatusNotFound,
		}
	}
	if err != nil {
		return db.LoginInfo{}, &models.Error{
			Msg:    fmt.Sprintf("failed to query raw login: %v", err),
			Status: http.StatusInternalServerError,
		}
	}

	return info, nil
}

func (r repository) GetUserInfo(_ context.Context, uuid string) (models.UserInfoRes, error) {
	query := `SELECT username, notifications_on,email FROM users WHERE uuid = ? LIMIT 1`

	var info models.UserInfoRes
	err := r.db.QueryRowx(query, uuid).StructScan(&info)
	if errors.Is(sql.ErrNoRows, err) {
		return models.UserInfoRes{}, &models.Error{
			Msg:    fmt.Sprintf("user with uuid: %s not found", uuid),
			Status: http.StatusNotFound,
		}
	}
	if err != nil {
		return models.UserInfoRes{}, &models.Error{
			Msg:    fmt.Sprintf("failed to query raw get user info: %v", err),
			Status: http.StatusInternalServerError,
		}
	}

	return info, nil
}

func (r repository) GetInfoForRent(_ context.Context, uuid string) (models.InfoForRentRes, error) {
	query := `SELECT passport_number,phone_number FROM users WHERE uuid = ? LIMIT 1`

	var info models.InfoForRentRes
	err := r.db.QueryRowx(query, uuid).StructScan(&info)
	if errors.Is(sql.ErrNoRows, err) {
		return models.InfoForRentRes{}, &models.Error{
			Msg:    fmt.Sprintf("user with uuid: %s not found", uuid),
			Status: http.StatusNotFound,
		}
	}
	if err != nil {
		return models.InfoForRentRes{}, &models.Error{
			Msg:    fmt.Sprintf("failed to query raw get user info for rent: %v", err),
			Status: http.StatusInternalServerError,
		}
	}

	return info, nil
}

func (r repository) SwitchNotificationsStatus(_ context.Context, uuid string) error {
	query := `UPDATE users SET notifications_on = !notifications_on WHERE uuid = ?`

	err := r.db.QueryRowx(query, uuid).Err()
	if errors.Is(sql.ErrNoRows, err) {
		return &models.Error{
			Msg:    fmt.Sprintf("user with uuid: %s not found", uuid),
			Status: http.StatusNotFound,
		}
	}
	if err != nil {
		return &models.Error{
			Msg:    fmt.Sprintf("failed to query raw switch notifications status: %v", err),
			Status: http.StatusInternalServerError,
		}
	}

	return nil
}

func (r repository) GetUserByUUID(_ context.Context, uuid string) (db.EmailNotificationsInfo, error) {
	query := `SELECT username, email FROM users WHERE uuid = ? LIMIT 1`

	var info db.EmailNotificationsInfo
	err := r.db.QueryRowx(query, uuid).StructScan(&info)
	if errors.Is(sql.ErrNoRows, err) {
		return db.EmailNotificationsInfo{}, &models.Error{
			Msg:    fmt.Sprintf("user with uuid: %s not found", uuid),
			Status: http.StatusNotFound,
		}
	}
	if err != nil {
		return db.EmailNotificationsInfo{}, &models.Error{
			Msg:    fmt.Sprintf("failed to query raw get user by uuid: %v", err),
			Status: http.StatusInternalServerError,
		}
	}

	return info, nil
}

const ROLE_USER = "user"

func (r repository) Register(_ context.Context, req models.RegisterReq) error {
	query := `INSERT INTO users (uuid,username,password,email,role ,passport_number,payment_source,phone_number)
				VALUES (?,?,?,?,?,?,?, ?)`

	err := r.db.QueryRowx(query, req.UUID, req.Username, req.Password, req.Email, ROLE_USER, req.PassportNumber, req.PaymentSource, req.PhoneNumber).Err()
	if err != nil {
		return &models.Error{
			Msg:    fmt.Sprintf("failed to query raw register: %v", err),
			Status: http.StatusInternalServerError,
		}
	}

	return nil
}
