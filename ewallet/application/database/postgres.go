package database

import (
	"gorm.io/gorm"
	"simple-go/application/config"
	"simple-go/application/entity"
	"simple-go/pkg/db"
	"time"
)

func ConnectPostgres() (*gorm.DB, error) {
	stageApp := config.GetString(config.CFG_APP_ENV, "")
	postgres := db.NewGormPostgres(
		config.GetString(config.CFG_POSTGRES_HOST, ""),
		config.GetString(config.CFG_POSTGRES_PORT, ""),
		config.GetString(config.CFG_POSTGRES_USER, ""),
		config.GetString(config.CFG_POSTGRES_PASS, ""),
		config.GetString(config.CFG_POSTGRES_DB, ""),
		config.GetString(config.CFG_POSTGRES_SSLMODE, ""),
	)

	err := postgres.Connect()
	if err != nil {
		return nil, err
	}

	err = postgres.SetConnectionPool(
		config.GetInt(config.CFG_POSTGRES_MAX_OPEN_CONNS, 0),
		config.GetInt(config.CFG_POSTGRES_MAX_IDLE_CONNS, 0),
		time.Duration(config.GetInt(config.CFG_POSTGRES_LIFETIME_IDLE_CONNS, 0)),
		time.Duration(config.GetInt(config.CFG_POSTGRES_LIFETIME_OPEN_CONNS, 0)),
	)
	if err != nil {
		return nil, err
	}

	db := postgres.(*db.GormPostgresDB).DB
	if stageApp == "dev" {
		db.AutoMigrate(
			&entity.Customer{},
			&entity.CustomerAuth{},
			&entity.Wallet{},
			&entity.WalletTransaction{},
		)
	}

	return db, nil
}
