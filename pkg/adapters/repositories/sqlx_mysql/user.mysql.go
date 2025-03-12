package sqlx_mysql

import (
	"go-clean-api/pkg/adapters/db"
	"go-clean-api/pkg/adapters/repositories/sqlx_mysql/models"
	"go-clean-api/pkg/domain/repositories"

	"github.com/jmoiron/sqlx"
)

// User is an implementation of the UserRepository interface
type User struct {
	db *sqlx.DB
}

// NewUser creates a new UserMysqlRepository
func NewUser(db *db.MySQL) *User {
	return &User{db: db.DB}
}

// GetByEmail returns user ID and password from the email
func (u *User) GetByEmail(req repositories.GetByEmailRequest) (repositories.GetByEmailResponse, error) {
	var model models.GetByEmail
	row := u.db.QueryRowx(`
		SELECT id, password
		FROM users
		WHERE email = ?
			AND deleted_at IS NULL
		LIMIT 1`,
		req.Email,
	)
	if err := row.StructScan(&model); err != nil {
		return repositories.GetByEmailResponse{}, repositories.ErrUserNotFound
	}

	response, err := model.ToRepository()
	if err != nil {
		return repositories.GetByEmailResponse{}, repositories.ErrConvertFromModel
	}

	return response, nil
}

func (u *User) Create(req repositories.CreateUserRequest) (repositories.CreateUserResponse, error) {
	// TODO: implement

	return repositories.CreateUserResponse{}, nil
}
