package model

import (
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type dbHandler struct {
	db *gorm.DB
}

func newSqliteHandler(path string) (DBHandler, error) {
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&Device{})
	db.AutoMigrate(&Controller{})
	db.AutoMigrate(&Service{})

	return &dbHandler{db}, nil
}

func newPostgresqlHandler(path string) (DBHandler, error) {
	dsn := "host=localhost user=user password=user_password dbname=godopudb port=5432 sslmode=disable TimeZone=Asia/Seoul"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&Device{})
	db.AutoMigrate(&Controller{})
	db.AutoMigrate(&Service{})

	return &dbHandler{db}, nil
}
