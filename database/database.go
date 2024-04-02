package database

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

var ErrDriver = errors.New("not support driver")

// NewDB 是通过配置项，创建一个数据库连接。
func NewDB(opts ...Option) (*gorm.DB, error) {
	c := newOptions(opts...)
	dialector, ok := opens[c.driver]
	if !ok {
		return nil, ErrDriver
	}
	db, err := gorm.Open(dialector(c.dsn), c.gormConfig)
	if err != nil {
		return nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return db, err
	}
	if c.connMaxIdleTime > 0 {
		sqlDB.SetConnMaxIdleTime(time.Duration(c.connMaxIdleTime) * time.Second)
	}
	if c.connMaxLifetime > 0 {
		sqlDB.SetConnMaxLifetime(time.Duration(c.connMaxLifetime) * time.Second)
	}
	if c.maxOpenConns > 0 {
		sqlDB.SetMaxOpenConns(c.maxOpenConns)
	}
	if c.maxIdleConns > 0 {
		sqlDB.SetMaxIdleConns(c.maxIdleConns)
	}
	return db, err
}
