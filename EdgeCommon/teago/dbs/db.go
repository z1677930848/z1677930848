package dbs

import (
	"database/sql"
	"errors"

	"github.com/iwind/TeaGo/maps"
)

// DB wraps a database handle (stubbed).
type DB struct {
	config *DBConfig
	sqlDB  *sql.DB
}

// Config returns underlying config.
func (db *DB) Config() (*DBConfig, error) {
	if db == nil {
		return nil, errors.New("nil db")
	}
	return db.config, nil
}

// Name returns identifier of DB.
func (db *DB) Name() string {
	if cfg, _ := db.Config(); cfg != nil {
		return cfg.Dsn
	}
	return ""
}

// Tx represents a transaction (stubbed).
type Tx struct {
	db *DB
}

// DBRaw is a simple wrapper providing Ping.
type DBRaw struct {
	db *DB
}

// NewInstanceFromConfig creates a DB instance from config.
func NewInstanceFromConfig(config *DBConfig) (*DB, error) {
	if config == nil {
		return nil, errors.New("nil db config")
	}
	return &DB{config: config}, nil
}

// Default returns a default DB using GlobalConfig and Tea.Env if available.
func Default() (*DB, error) {
	if globalConfig == nil || len(globalConfig.DBs) == 0 {
		return nil, errors.New("no db config")
	}
	// pick the first config
	for _, cfg := range globalConfig.DBs {
		return NewInstanceFromConfig(cfg)
	}
	return nil, errors.New("no db config")
}

// Close closes underlying db.
func (db *DB) Close() error {
	if db == nil || db.sqlDB == nil {
		return nil
	}
	return db.sqlDB.Close()
}

// SetMaxOpenConns sets max open connections.
func (db *DB) SetMaxOpenConns(n int) {
	if db != nil && db.sqlDB != nil {
		db.sqlDB.SetMaxOpenConns(n)
	}
}

// SetMaxIdleConns sets max idle connections.
func (db *DB) SetMaxIdleConns(n int) {
	if db != nil && db.sqlDB != nil {
		db.sqlDB.SetMaxIdleConns(n)
	}
}

// Raw returns a raw wrapper.
func (db *DB) Raw() *DBRaw {
	return &DBRaw{db: db}
}

// Ping checks connectivity (stubbed).
func (raw *DBRaw) Ping() error {
	if raw == nil || raw.db == nil || raw.db.sqlDB == nil {
		return nil
	}
	return raw.db.sqlDB.Ping()
}

// PingContext is a context-aware ping (stub).
func (raw *DBRaw) PingContext(ctx any) error {
	return raw.Ping()
}

// SetMaxOpenConns proxies to underlying DB.
func (raw *DBRaw) SetMaxOpenConns(n int) {
	if raw != nil && raw.db != nil {
		raw.db.SetMaxOpenConns(n)
	}
}

// Exec executes a statement (stubbed).
func (db *DB) Exec(query string, args ...any) (sql.Result, error) {
	if db != nil && db.sqlDB != nil {
		return db.sqlDB.Exec(query, args...)
	}
	return nil, nil
}

// FindOne executes query and returns first row as maps.Map (stubbed).
func (db *DB) FindOne(query string, args ...any) (maps.Map, error) {
	return maps.Map{}, nil
}

// FindCol returns a column value from query.
func (db *DB) FindCol(index int, query string, args ...any) (any, error) {
	return nil, nil
}

// FindOnes returns multiple rows from a query (stubbed).
func (db *DB) FindOnes(_ string, _ ...any) ([]maps.Map, int64, error) {
	return []maps.Map{}, 0, nil
}

// TableNames lists all table names in current database (stubbed).
func (db *DB) TableNames() ([]string, error) {
	return []string{}, nil
}

// FindFullTable returns full table schema info (stubbed).
func (db *DB) FindFullTable(_ string) (*Table, error) {
	return &Table{}, nil
}

// FindTable returns schema for a table (stubbed).
func (db *DB) FindTable(_ string) (*Table, error) {
	return &Table{}, nil
}

// Prepare prepares a statement (stubbed).
func (db *DB) Prepare(query string) (*Stmt, error) {
	return &Stmt{db: db, query: query}, nil
}

// StmtManager returns a prepared statement manager stub.
func (db *DB) StmtManager() *StmtManager {
	return &StmtManager{}
}

// FindPreparedOnes executes a prepared statement and returns rows and column names.
func (db *DB) FindPreparedOnes(_ string, _ ...any) ([]maps.Map, []string, error) {
	return []maps.Map{}, []string{}, nil
}

// Begin starts a transaction (stubbed).
func (db *DB) Begin() (*Tx, error) {
	return &Tx{db: db}, nil
}

// RunTx executes fn within a transaction (stub).
func (db *DB) RunTx(fn func(tx *Tx) error) error {
	if fn == nil {
		return nil
	}
	tx, _ := db.Begin()
	return fn(tx)
}

// Commit commits transaction (stubbed).
func (tx *Tx) Commit() error {
	return nil
}

// Rollback rollbacks transaction (stubbed).
func (tx *Tx) Rollback() error {
	return nil
}

// FindCol delegates to DB FindCol when available.
func (tx *Tx) FindCol(index int, query string, args ...any) (any, error) {
	if tx != nil && tx.db != nil {
		return tx.db.FindCol(index, query, args...)
	}
	return nil, nil
}
