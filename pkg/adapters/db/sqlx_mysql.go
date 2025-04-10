package db

import (
	"go-clean-api/pkg"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// SqlxMySQL is a struct that contains the database connection
type SqlxMySQL struct {
	DB     *sqlx.DB
	config *pkg.Config
}

// NewSqlxMySQL creates a new MySQL database connection
func NewSqlxMySQL(config *pkg.Config) (*SqlxMySQL, error) {
	dsn, err := config.Database.DSN()
	if err != nil {
		return nil, err
	}

	db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	db.SetConnMaxIdleTime(config.Database.ConnMaxIdleTime)
	db.SetConnMaxLifetime(config.Database.ConnMaxLifetime)
	db.SetMaxOpenConns(config.Database.MaxOpenConns)
	db.SetMaxIdleConns(config.Database.MaxIdleConns)

	return &SqlxMySQL{
		DB:     db,
		config: config,
	}, nil
}

func (m *SqlxMySQL) DSN() (string, error) {
	return m.config.Database.DSN()
}

func (m *SqlxMySQL) Database(d string) {
	m.config.Database.Database = d
}
