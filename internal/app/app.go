package app

import (
	"fmt"
	"go-clean-api/pkg"
	"go-clean-api/pkg/adapters/db"
	"go-clean-api/pkg/adapters/repositories/gorm_mysql"
	"go-clean-api/pkg/domain/usecases"
	"go-clean-api/pkg/infrastructure/auth"
	"go-clean-api/pkg/infrastructure/logger"
)

// Dependencies holds all wired dependencies for the application.
type Dependencies struct {
	Config      pkg.Config
	DB          db.DB
	Logger      logger.CustomLogger
	UserUseCase usecases.User
}

// NewDependencies creates and wires all application dependencies.
func NewDependencies(config pkg.Config, database db.DB, l logger.CustomLogger) (*Dependencies, error) {
	gormDB, ok := database.(*db.GormMySQL)
	if !ok {
		return nil, fmt.Errorf("database is not of type *db.GormMySQL")
	}

	userRepo := gorm_mysql.NewUser(gormDB)
	tokenGen := auth.NewJWTTokenGenerator(config.JWT)
	userUseCase := usecases.NewUser(userRepo, tokenGen)

	return &Dependencies{
		Config:      config,
		DB:          database,
		Logger:      l,
		UserUseCase: userUseCase,
	}, nil
}
