package db

import (
	"fmt"
	"sync"

	"github.com/iwind/TeaGo/Tea"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	connMu   sync.Mutex
	connPool = map[string]*gorm.DB{}
	dbConfig *Config
)

// Init loads config once.
func Init() error {
	if dbConfig != nil {
		return nil
	}
	cfg, err := LoadConfig()
	if err != nil {
		return err
	}
	dbConfig = cfg
	return nil
}

// DB returns a gorm.DB for the given env (Tea.Env by default).
func DB(env ...string) (*gorm.DB, error) {
	if err := Init(); err != nil {
		return nil, err
	}
	e := Tea.Env
	if len(env) > 0 && env[0] != "" {
		e = env[0]
	}
	cfg, ok := dbConfig.DBs[e]
	if !ok {
		return nil, fmt.Errorf("db config for env %q not found", e)
	}

	connMu.Lock()
	defer connMu.Unlock()
	if db, ok := connPool[e]; ok {
		return db, nil
	}

	var dial gorm.Dialector
	switch cfg.Driver {
	case "mysql":
		dial = mysql.Open(cfg.DSN)
	default:
		return nil, fmt.Errorf("unsupported driver %q", cfg.Driver)
	}

	db, err := gorm.Open(dial, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		return nil, err
	}
	connPool[e] = db
	return db, nil
}
