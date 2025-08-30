package mysql

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/thanhfphan/kart-challenge/config"

	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitConn(cfg *config.DB) (*gorm.DB, *sql.DB, error) {
	db, err := gorm.Open(mysql.Open(cfg.ConnectionURL), &gorm.Config{})
	if err != nil {
		return nil, nil, err
	}

	if err := db.Use(otelgorm.NewPlugin()); err != nil {
		return nil, nil, fmt.Errorf("new otel-gorm plugin got err=%w", err)
	}

	dbConfig, err := db.DB()
	if err != nil {
		return nil, nil, err
	}

	dbConfig.SetMaxOpenConns(cfg.MaxOpenConnNumber)
	dbConfig.SetMaxIdleConns(cfg.MaxIdleConnNumber)
	dbConfig.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifeTimeSeconds) * time.Second)
	if err = dbConfig.Ping(); err != nil {
		return nil, nil, err
	}

	return db, dbConfig, nil
}
