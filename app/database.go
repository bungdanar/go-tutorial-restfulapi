package app

import (
	"database/sql"
	"time"
	"tutorial-restfulapi/helper"
)

func NewDB() *sql.DB {
	db, err := sql.Open("mysql", "root:brightshield!23@tcp(localhost:3306)/go-tutorial-restfulapi")
	helper.PanicIfError(err)

	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(20)
	db.SetConnMaxIdleTime(10 * time.Minute)
	db.SetConnMaxLifetime(60 * time.Minute)

	return db
}
