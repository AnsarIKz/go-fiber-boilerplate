package repository

import (
	"nodabackend/internal/auth/domain"

	"gorm.io/gorm"
)

// AuthRepository реализация AuthRepository для PostgreSQL
type AuthRepository struct {
	db *gorm.DB
}

// NewAuthRepository создает новый репозиторий с dependency injection
func NewAuthRepository(db *gorm.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

// CreateUser создает нового пользователя
func (r *AuthRepository) CreateUser(user *domain.User) error {
	return r.db.Create(user).Error
}

// GetUserByPhone получает пользователя по телефону
func (r *AuthRepository) GetUserByPhone(phone string) (*domain.User, error) {
	var user domain.User
	err := r.db.Where("phone = ?", phone).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByID получает пользователя по ID
func (r *AuthRepository) GetUserByID(id uint) (*domain.User, error) {
	var user domain.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
