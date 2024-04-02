package database

import (
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

const (
	DEFAULT_DB_DRIVER = "sqlite3"
	DEFAULT_DB_DSN    = ":memory:"
)

type Options struct {
	driver          string
	dsn             string
	connMaxIdleTime int
	connMaxLifetime int
	maxIdleConns    int
	maxOpenConns    int
	gormConfig      *gorm.Config
}

type Option func(o *Options)

func newOptions(options ...Option) Options {
	opts := Options{
		driver:          DEFAULT_DB_DRIVER,
		dsn:             DEFAULT_DB_DSN,
		connMaxIdleTime: 60 * 5,
		connMaxLifetime: 60 * 60,
		maxIdleConns:    2,
		maxOpenConns:    25,
		gormConfig: &gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true,
			},
		},
	}

	for _, o := range options {
		o(&opts)
	}
	return opts
}

func WithDriver(driver string) Option {
	return func(o *Options) {
		o.driver = driver
	}
}

func WithDSN(dsn string) Option {
	return func(o *Options) {
		o.dsn = dsn
	}
}

func WithConnMaxIdleTime(t int) Option {
	return func(o *Options) {
		o.connMaxIdleTime = t
	}
}

func WithConnMaxLifetime(t int) Option {
	return func(o *Options) {
		o.connMaxLifetime = t
	}
}

func WithMaxIdleConns(c int) Option {
	return func(o *Options) {
		o.maxIdleConns = c
	}
}

func WithMaxOpenConns(c int) Option {
	return func(o *Options) {
		o.maxOpenConns = c
	}
}

func WithGormConfig(c *gorm.Config) Option {
	return func(o *Options) {
		o.gormConfig = c
	}
}
