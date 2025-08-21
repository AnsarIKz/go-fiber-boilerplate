package domain

import "time"

// User представляет пользователя в системе
type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Phone     string    `json:"phone" gorm:"uniqueIndex"`
	Password  string    `json:"-" gorm:"column:password"`  // пароль не возвращается в JSON
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// AuthRepository интерфейс для работы с аутентификацией
type AuthRepository interface {
	CreateUser(user *User) error
	GetUserByPhone(phone string) (*User, error)
	GetUserByID(id uint) (*User, error)
}
