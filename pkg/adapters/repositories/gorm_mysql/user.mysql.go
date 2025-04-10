package gorm_mysql

import (
	"go-clean-api/pkg/adapters/db"

	"gorm.io/gorm"
)

// User is an implementation of the UserRepository interface
type User struct {
	db *gorm.DB
}

// NewUser creates a new UserMysqlRepository
func NewUser(db *db.GormMySQL) *User {
	return &User{db: db.DB}
}
