package domain

import "time"

// Role представляет роль пользователя в системе
type Role string

const (
	AdminRole Role = "admin"
	UserRole  Role = "user"
)

// User представляет пользователя в системе
type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Phone     string    `json:"phone" gorm:"uniqueIndex"`
	Password  string    `json:"-" gorm:"column:password"`  // пароль не возвращается в JSON
	Name      string    `json:"name"`
	Role      Role      `json:"role" gorm:"type:varchar(20);default:'user'"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UserRepository интерфейс для работы с пользователями
type UserRepository interface {
	CreateUser(user *User) error
	GetUserByPhone(phone string) (*User, error)
	GetUserByID(id uint) (*User, error)
}
