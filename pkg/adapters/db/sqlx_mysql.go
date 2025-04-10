package db

import (
	"go-clean-api/pkg"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// MySQL is a struct that contains the database connection
type MySQL struct {
	DB     *sqlx.DB
	config *pkg.Config
}

// NewMySQL creates a new MySQL database connection
func NewMySQL(config *pkg.Config) (*MySQL, error) {
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

	return &MySQL{
		DB:     db,
		config: config,
	}, nil
}

func (m *MySQL) DSN() (string, error) {
	return m.config.Database.DSN()
}

func (m *MySQL) Database(d string) {
	m.config.Database.Database = d
}
