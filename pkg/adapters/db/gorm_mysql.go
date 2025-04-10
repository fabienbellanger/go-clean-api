package db

import (
	"go-clean-api/pkg"
	"io"
	"log"
	"os"
	"path"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// GormMySQL is a struct that contains the database connection using Gorm ORM
type GormMySQL struct {
	DB     *gorm.DB
	config *pkg.Config
}

func NewGormMySQL(config *pkg.Config) (*GormMySQL, error) {
	dsn, err := config.Database.DSN()
	if err != nil {
		return nil, err
	}

	if config.Gorm.SlowThreshold == 0 {
		config.Gorm.SlowThreshold = DefaultSlowThreshold
	}

	// GORM logger configuration
	env := config.AppEnv
	level := getGormLogLevel(config.Gorm.LogLevel, env)
	output, err := getGormLogOutput(config.Gorm.LogOutput, config.Gorm.LogFileName, env)
	if err != nil {
		return nil, err
	}

	// Logger
	// TODO: Add a custom logger for GORM like https://www.soberkoder.com/go-gorm-logging/
	// Or try something like this: https://github.com/moul/zapgorm2
	customLogger := logger.New(
		log.New(output, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             config.Gorm.SlowThreshold, // Slow SQL threshold (Default: 200ms)
			LogLevel:                  level,                     // Log level (Silent, Error, Warn, Info) (Default: Warn)
			IgnoreRecordNotFoundError: true,                      // Ignore ErrRecordNotFound error for logger (Default: false)
			Colorful:                  true,                      // Disable color (Default: true)
		},
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: customLogger,
	})
	if err != nil {
		return nil, err
	}

	// Options
	// -------
	db.Set("gorm:table_options", "ENGINE=InnoDB")

	// Connection Pool
	// ---------------
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetConnMaxIdleTime(config.Database.ConnMaxIdleTime)
	sqlDB.SetConnMaxLifetime(config.Database.ConnMaxLifetime)
	sqlDB.SetMaxOpenConns(config.Database.MaxOpenConns)
	sqlDB.SetMaxIdleConns(config.Database.MaxIdleConns)

	return &GormMySQL{
		DB:     db,
		config: config,
	}, nil
}

func (m *GormMySQL) DSN() (string, error) {
	return m.config.Database.DSN()
}

func (m *GormMySQL) Database(d string) {
	m.config.Database.Database = d
}

// getGormLogLevel returns the log level for GORM.
// If APP_ENV is development, the default log level is info,
// warn in other case.
func getGormLogLevel(level, env string) logger.LogLevel {
	switch level {
	case "silent":
		return logger.Silent
	case "info":
		return logger.Info
	case "warn":
		return logger.Warn
	case "error":
		return logger.Error
	default:
		if env == "development" {
			return logger.Warn
		}
		return logger.Error
	}
}

// getGormLogOutput returns GORM log output.
// The default value is os.Stdout.
// In development mode, the ouput is set to os.Stdout.
func getGormLogOutput(output, filePath, env string) (file io.Writer, err error) {
	if env == "development" {
		return os.Stdout, nil
	}

	switch output {
	case "file":
		f, err := os.OpenFile(path.Clean(filePath), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			return nil, err
		}
		return f, nil
	default:
		return os.Stdout, nil
	}
}

// GormPaginate creates a GORM scope to paginate queries.
func GormPaginate(p, l int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		offset, limit := PaginateValues(p, l)

		return db.Offset(offset).Limit(limit)
	}
}

// GormOrder creates a GORM scope to sort query attributes.
// Example: "+created_at,-id" will produce "ORDER BY created_at ASC, id DESC".
func GormOrder(list string, prefixes ...string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		values := orderValues(list, prefixes...)

		for s := range values {
			db.Order(s)
		}

		return db
	}
}
