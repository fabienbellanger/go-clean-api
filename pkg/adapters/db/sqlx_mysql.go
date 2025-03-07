package db

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// MySQL is a struct that contains the database connection
type MySQL struct {
	DB     *sqlx.DB
	config *Config
}

// NewMySQL creates a new MySQL database connection
func NewMySQL(config *Config) (*MySQL, error) {
	dsn, err := config.dsn()
	if err != nil {
		return nil, err
	}

	db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	db.SetConnMaxIdleTime(config.ConfigDatabase.ConnMaxIdleTime)
	db.SetConnMaxLifetime(config.ConnMaxLifetime)
	db.SetMaxOpenConns(config.MaxOpenConns)
	db.SetMaxIdleConns(config.MaxIdleConns)

	return &MySQL{
		DB:     db,
		config: config,
	}, nil
}

func (m *MySQL) DSN() (string, error) {
	return m.config.dsn()
}

func (m *MySQL) Database(d string) {
	m.config.Database = d
}
