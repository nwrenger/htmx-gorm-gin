package db

import (
	"database/sql"
	"fmt"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func Open() error {
	connString := "htmx-gorm-gin.db"

	// open SQLite database
	sqlDB, err := sql.Open("sqlite3", connString)
	if err != nil {
		return err
	}
	sqlDB.SetConnMaxLifetime(10 * time.Minute)
	sqlDB.SetMaxOpenConns(50)

	// hand over to gorm with SQLite driver
	db, err = gorm.Open(sqlite.Open(connString), &gorm.Config{})
	if err != nil {
		return err
	}

	// todo logging
	fmt.Println("successfully connected to SQLite database")

	// specify and auto-migrate tables
	// the order here is important
	return db.AutoMigrate(
		&Book{},
		&Login{},
	)
}
