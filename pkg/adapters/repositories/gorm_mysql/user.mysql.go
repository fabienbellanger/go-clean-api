package gorm_mysql

import (
	"fmt"
	"go-clean-api/pkg/adapters/db"
	"go-clean-api/pkg/adapters/models"
	"go-clean-api/pkg/domain/entities"
	"go-clean-api/pkg/domain/repositories"

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

func (u *User) GetByEmail(req repositories.GetByEmailRequest) (res repositories.GetByEmailResponse, err error) {
	var model models.GetUserByEmail
	result := u.db.Raw(`
		SELECT id, password
		FROM users
		WHERE email = ?
			AND deleted_at IS NULL
		LIMIT 1`, req.Email.Value()).Scan(&model)
	if result.Error != nil {
		return res, fmt.Errorf("[user_gorm_mysql:GetByEmail %w: %s]", repositories.ErrUserNotFound, result.Error)
	} else if result.RowsAffected == 0 {
		return res, fmt.Errorf("[user_gorm_mysql:GetByEmail %w]", repositories.ErrUserNotFound)
	}

	res, err = model.Repository()

	if err != nil {
		return res, fmt.Errorf("[user_gorm_mysql:GetByEmail %w: %s]", repositories.ErrGettingUser, err)
	}

	return res, nil
}

func (u *User) GetByID(req repositories.GetByIDRequest) (res repositories.GetByIDResponse, err error) {
	var model models.User
	if result := u.db.Raw(`
		SELECT id, email, lastname, firstname, created_at, updated_at, deleted_at
		FROM users
		WHERE id = ?
			AND deleted_at IS NULL
		LIMIT 1`, req.ID.Value()).Scan(&model); result.Error != nil {
		return res, fmt.Errorf("[user_gorm_mysql:GetByID %w: %s]", repositories.ErrUserNotFound, result.Error)
	}
	user, err := model.Entity()

	if err != nil {
		return res, fmt.Errorf("[user_gorm_mysql:GetByID %w: %s]", repositories.ErrGettingUser, err)
	}

	res.User = user

	return
}

func (u *User) CountAll(req repositories.CountAllRequest) (repositories.CountAllResponse, error) {
	q := `
		SELECT COUNT(id) AS total
		FROM users`

	if req.Deleted {
		q += " WHERE deleted_at IS NOT NULL"
	} else {
		q += " WHERE deleted_at IS NULL"
	}

	var count int64
	row := u.db.Raw(q)
	if result := row.Scan(&count); result.Error != nil {
		return repositories.CountAllResponse{}, fmt.Errorf("[user_gorm_mysql:CountAll %w: %s]", repositories.ErrCountingUsers, result.Error)
	}

	return repositories.CountAllResponse{Total: count}, nil
}

func (u *User) GetAll(req repositories.GetAllRequest) (res repositories.GetAllResponse, err error) {
	q := u.db.Scopes(db.GormPaginate(req.Pagination.Page(), req.Pagination.Size()))
	if req.Deleted {
		q = q.Where("deleted_at IS NOT NULL")
	} else {
		q = q.Where("deleted_at IS NULL")
	}

	var users []models.User
	if result := q.Find(&users); result.Error != nil {
		return res, fmt.Errorf("[user_gorm_mysql:GetAll %w: %s]", repositories.ErrGettingUsers, result.Error)
	}

	usersEntity := make([]entities.User, 0, len(users))
	for _, user := range users {
		userEntity, err := user.Entity()
		if err != nil {
			return res, fmt.Errorf("[user_gorm_mysql:GetAll %w: %s]", repositories.ErrGettingUsers, err)
		}
		usersEntity = append(usersEntity, userEntity)
	}
	res.Users = usersEntity

	return
}

func (u *User) Create(req repositories.CreateUserRequest) (res repositories.CreateUserResponse, err error) {
	result := u.db.Exec(`
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
	if result.Error != nil {
		return res, fmt.Errorf("[user_gorm_mysql:Create %w: %s]", repositories.ErrCreatingUser, result.Error)
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

func (u *User) Delete(req repositories.DeleteRestoreRequest) (res repositories.DeleteRestoreResponse, err error) {
	result := u.db.Exec(`
		UPDATE users
		SET deleted_at = NOW()
		WHERE id = ?
			AND deleted_at IS NULL`,
		req.ID.String(),
	)
	if result.Error != nil {
		return res, fmt.Errorf("[user_sqlx_mysql:Delete %w: %s]", repositories.ErrDatabase, result.Error)
	}

	if result.RowsAffected == 0 {
		return repositories.DeleteRestoreResponse{}, fmt.Errorf("[user_sqlx_mysql:Delete %w]", repositories.ErrUserNotFound)
	}

	return
}

func (u *User) Restore(req repositories.DeleteRestoreRequest) (res repositories.DeleteRestoreResponse, err error) {
	result := u.db.Exec(`
		UPDATE users
		SET deleted_at = NULL
		WHERE id = ?
			AND deleted_at IS NOT NULL`,
		req.ID.String(),
	)
	if result.Error != nil {
		return res, fmt.Errorf("[user_sqlx_mysql:Restore %w: %s]", repositories.ErrDatabase, result.Error)
	}

	if result.RowsAffected == 0 {
		return repositories.DeleteRestoreResponse{}, fmt.Errorf("[user_sqlx_mysql:Restore %w]", repositories.ErrUserNotFound)
	}

	return
}
