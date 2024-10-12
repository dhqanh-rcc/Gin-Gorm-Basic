package util

import (
	"fmt"
	"log"
	"time"

	"github.com/quanganh247-qa/gorm-project/app/db"
	"github.com/quanganh247-qa/gorm-project/app/util/token"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Connection struct {
	Close func()
	// Db    *gorm.DB
}

func Init(config Config) (*Connection, error) {

	// Initialize JWT token maker
	_, err := token.NewJWTMaker(config.SymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("can't create token maker: %w", err)
	}

	sqldb, err := gorm.Open(postgres.Open(config.DSN), &gorm.Config{})
	if err != nil {
		log.Fatal("Error connecting to DB:", err)
	}

	sqlDB, err := sqldb.DB()
	sqlDB.SetMaxOpenConns(25)                 // Maximum number of open connections to the database
	sqlDB.SetMaxIdleConns(25)                 // Maximum number of connections in the idle connection pool
	sqlDB.SetConnMaxLifetime(5 * time.Minute) // Maximum amount of time a connection may be reused

	db.InitStore(sqldb)
	conn := &Connection{
		Close: func() {
			sqlDB.Close()
		},
		// Db: sqldb,
	}
	return conn, nil
}
