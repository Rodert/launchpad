package repository

import (
	"errors"
	"fmt"
	"launchpad/config"
	"launchpad/database"
	"launchpad/database/gormlog"
	olog "launchpad/logger"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var dbMap map[string]*gorm.DB

func init() {
	dbMap = make(map[string]*gorm.DB)
}

func InitPg(config config.Configuration) bool {
	for name, cfg := range config.PostgreCfg.Conf {
		if cfg == "" {
			continue
		}

		db, err := initPg(cfg, config.Log)
		if err == nil && db != nil {
			dbMap[name] = db
		} else {
			panic(err)
		}
	}
	return true
}

func initPg(path string, logCfg olog.Conf) (*gorm.DB, error) {
	var (
		db  *gorm.DB
		err error
	)
	// olog.Entry.Errorf("failed init db.[ path = %s]", path)
	fmt.Println(path)
	myLog, _ := olog.InitLog(olog.Conf{
		LogLevel: logCfg.LogLevel,
		LogPath:  logCfg.LogPath,
	})
	gormConfig := &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		Logger: gormlog.New(myLog.Logger, logger.Config{
			SlowThreshold:             time.Second,
			Colorful:                  false,
			IgnoreRecordNotFoundError: false,
			LogLevel:                  logger.Info,
		}),
	}
	opts := []database.Option{
		// 配置驱动，可选驱动到
		// database.WithDriver("sqlite3"),
		// database.WithDSN("test.db"),
		database.WithDriver("postgres"),
		database.WithDSN(path),
		database.WithMaxIdleConns(2),
		database.WithMaxOpenConns(10),
		database.WithConnMaxIdleTime(60 * 5),
		database.WithConnMaxLifetime(60 * 60),
		database.WithGormConfig(gormConfig),
	}
	for i := 0; i < 3; i++ {
		db, err = database.NewDB(opts...)
		if err == nil {
			break
		}
		time.Sleep(time.Second * 5)
	}
	if err != nil {
		return nil, err
	}
	// olog.Entry.Errorf("failed init db.[ path = %s err = %v]", path, err)
	fmt.Println(err)
	return db, err
}

func GetDB(name string) (*gorm.DB, error) {
	if db, ok := dbMap[name]; ok {
		return db, nil
	}

	errS := fmt.Sprintf("db not init: %s", name)

	return nil, errors.New(errS)
}
