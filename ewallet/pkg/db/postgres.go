package db

import (
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type GormPostgres interface {
	Connect() error
	SetConnectionPool(maxOpen int, maxIdle int, idleTime time.Duration, lifetime time.Duration) error
}

type GormPostgresDB struct {
	host   string
	port   string
	user   string
	pass   string
	dbname string
	ssl    string
	DB     *gorm.DB
}

func NewGormPostgres(host string, port string, user string, pass string, dbname string, ssl string) GormPostgres {
	return &GormPostgresDB{
		host:   host,
		port:   port,
		user:   user,
		pass:   pass,
		dbname: dbname,
		ssl:    ssl,
	}
}

func (g *GormPostgresDB) Connect() error {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		g.host, g.port, g.user, g.pass, g.dbname, g.ssl,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		TranslateError: true,
	})
	if err != nil {
		return err
	}
	g.DB = db

	return nil
}

func (g *GormPostgresDB) SetConnectionPool(maxOpen int, maxIdle int, idleTime time.Duration, lifetime time.Duration) error {
	db, err := g.DB.DB()
	if err != nil {
		return err
	}

	db.SetMaxOpenConns(maxOpen)
	db.SetMaxIdleConns(maxIdle)
	db.SetConnMaxIdleTime(idleTime * time.Second)
	db.SetConnMaxLifetime(lifetime * time.Second)

	return nil
}
