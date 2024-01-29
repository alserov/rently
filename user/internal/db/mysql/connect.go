package mysql

import "github.com/jmoiron/sqlx"

func MustConnect() *sqlx.DB {
	conn, err := sqlx.Open("mysql", "db_user:password@tcp(localhost:3306)/template1")
	if err != nil {
		panic("failed to open db: " + err.Error())
	}

	return conn
}
