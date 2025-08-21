package repository

import (
	"nodabackend/internal/auth/domain"

	"gorm.io/gorm"
)

// UserRepository реализация UserRepository для PostgreSQL
type UserRepository struct {
	db *gorm.DB
}


// NewUserRepository создает новый репозиторий с dependency injection
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// CreateUser создает нового пользователя
func (r *UserRepository) CreateUser(user *domain.User) error {
	return r.db.Create(user).Error
}

// GetUserByPhone получает пользователя по телефону
func (r *UserRepository) GetUserByPhone(phone string) (*domain.User, error) {
	var user domain.User
	err := r.db.Where("phone = ?", phone).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByID получает пользователя по ID
func (r *UserRepository) GetUserByID(id uint) (*domain.User, error) {
	var user domain.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
