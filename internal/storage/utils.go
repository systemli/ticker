package storage

import (
	"github.com/sirupsen/logrus"
	"github.com/systemli/ticker/internal/logger"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func OpenGormDB(dbType, dsn string, log *logrus.Logger) (*gorm.DB, error) {
	d := dialector(dbType, dsn)
	db, err := gorm.Open(d, &gorm.Config{
		Logger: logger.NewGormLogger(log),
	})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func dialector(dbType, dsn string) gorm.Dialector {
	switch dbType {
	case "sqlite":
		return sqlite.Open(dsn)
	case "mysql":
		return mysql.Open(dsn)
	case "postgres":
		return postgres.Open(dsn)
	default:
		return nil
	}
}
