package sqlx_mysql

import (
	"fmt"
	"go-clean-api/pkg/adapters/db"
	"go-clean-api/pkg/adapters/repositories/sqlx_mysql/models"
	"go-clean-api/pkg/domain/entities"
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
		req.Email.Value(),
	)
	if err := row.StructScan(&model); err != nil {
		return repositories.GetByEmailResponse{}, fmt.Errorf("[user_sqlx_mysql:GetByEmail] %w: (%v)", repositories.ErrUserNotFound, err)
	}

	response, err := model.Repository()
	if err != nil {
		return repositories.GetByEmailResponse{}, fmt.Errorf("[user_sqlx_mysql:GetByEmail] %w: (%v)", repositories.ErrConvertFromModel, err)
	}

	return response, nil
}

func (u *User) Create(req repositories.CreateUserRequest) (res repositories.CreateUserResponse, err error) {
	_, err = u.db.Exec(`
		INSERT INTO users (id, email, password, lastname, firstname, created_at, updated_at) 
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		req.ID.Value(),
		req.Email.Value(),
		req.Password.Value(),
		req.Lastname,
		req.Firstname,
		req.CreatedAt.SQL(),
		req.UpdatedAt.SQL(),
	)

	if err != nil {
		return
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

func (u *User) GetByID(req repositories.GetByIDRequest) (res repositories.GetByIDResponse, err error) {
	var model models.User
	row := u.db.QueryRowx(`
		SELECT id, email, lastname, firstname, created_at, updated_at, deleted_at
		FROM users
		WHERE id = ?
			AND deleted_at IS NULL
		LIMIT 1`,
		req.ID.String(),
	)
	if err = row.StructScan(&model); err != nil {
		return repositories.GetByIDResponse{}, fmt.Errorf("[user_sqlx_mysql:GetByID] %w: (%v)", repositories.ErrUserNotFound, err)
	}

	user, err := model.Entity()
	if err != nil {
		return repositories.GetByIDResponse{}, fmt.Errorf("[user_sqlx_mysql:GetByID] %w: (%v)", repositories.ErrConvertFromModel, err)
	}

	res.User = user

	return
}
