package mysql

import (
	"context"
	"github.com/alserov/rently/user/internal/db/migrations"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"os"
	"time"
)

func MustConnect(dsn string) *sqlx.DB {
	conn, err := sqlx.Open("mysql", dsn+"?multiStatements=true")
	if err != nil {
		panic("failed to open db: " + err.Error())
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	if err = conn.PingContext(ctx); err != nil {
		panic("failed to ping db: " + err.Error())
	}

	migrations.Migrate(conn)

	initAdmin(conn)

	return conn
}

func initAdmin(conn *sqlx.DB) {
	if err := godotenv.Load(".env"); err != nil {
		panic("failed to find .env file: " + err.Error())
	}

	adminPassword := os.Getenv("ADMIN_PASSWORD")
	panicIfEmpty(adminPassword, "admin password")
	adminEmail := os.Getenv("ADMIN_EMAIL")
	panicIfEmpty(adminEmail, "admin email")

	query := `INSERT INTO users (uuid,username,password,email,role ,passport_number,payment_source,phone_number)
				VALUES (?,?,?,?,?,?,?,?)`

	b, err := bcrypt.GenerateFromPassword([]byte(adminPassword), bcrypt.DefaultCost)
	if err != nil {
		panic("failed to hash password")
	}

	if err = conn.QueryRowx(query, uuid.New().String(), "ADMIN", string(b), adminEmail, "admin", "0", "0", "0").Err(); err != nil {
		panic("failed to init admin: " + err.Error())
	}
}

func panicIfEmpty(s string, msg string) {
	if s == "" {
		panic("can not be empty: " + msg)
	}
}
