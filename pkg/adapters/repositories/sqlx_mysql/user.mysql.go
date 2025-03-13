package sqlx_mysql

import (
	"fmt"
	"go-clean-api/pkg/adapters/db"
	"go-clean-api/pkg/adapters/repositories/sqlx_mysql/models"
	"go-clean-api/pkg/domain/entities"
	"go-clean-api/pkg/domain/repositories"
	"go-clean-api/utils"

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
		req.Email.Value(),
	)
	if err := row.StructScan(&model); err != nil {
		return repositories.GetByEmailResponse{}, fmt.Errorf("%w: [%v]", repositories.ErrUserNotFound, err)
	}

	response, err := model.ToRepository()
	if err != nil {
		return repositories.GetByEmailResponse{}, repositories.ErrConvertFromModel
	}

	return response, nil
}

func (u *User) Create(req repositories.CreateUserRequest) (repositories.CreateUserResponse, error) {
	_, err := u.db.Exec(`
		INSERT INTO users (id, email, password, lastname, firstname, created_at, updated_at) 
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		req.ID.Value(),
		req.Email.Value(),
		req.Password.Value(),
		req.Lastname,
		req.Firstname,
		utils.FormatToSqlDateTime(req.CreatedAt),
		utils.FormatToSqlDateTime(req.UpdatedAt),
	)

	if err != nil {
		return repositories.CreateUserResponse{}, err
	}

	return repositories.CreateUserResponse{
		User: entities.User{
			ID:        req.ID,
			Email:     req.Email,
			Password:  req.Password,
			Lastname:  req.Lastname,
			Firstname: req.Firstname,
			CreatedAt: req.CreatedAt,
			UpdatedAt: req.UpdatedAt,
		},
	}, nil
}
