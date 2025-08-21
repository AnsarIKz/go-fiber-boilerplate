package usecase

import (
	"errors"
	"nodabackend/internal/auth/domain"
	"nodabackend/pkg/jwthelper"
	"regexp"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// AuthResponse ответ при успешной аутентификации
type AuthResponse struct {
	User  *domain.User `json:"user"`
	Token string       `json:"token"`
}

// AuthUseCase бизнес-логика аутентификации
type AuthUseCase struct {
	userRepo  domain.UserRepository
	jwtHelper *jwthelper.JWTHelper
}

// NewAuthUseCase создает новый usecase
func NewAuthUseCase(repo domain.UserRepository, jwtHelper *jwthelper.JWTHelper) *AuthUseCase {
	return &AuthUseCase{
		userRepo:  repo,
		jwtHelper: jwtHelper,
	}
}

// RegisterUser регистрирует нового пользователя
func (uc *AuthUseCase) RegisterUser(phone, password, name string) (*AuthResponse, error) {
	// Валидация входных данных
	if err := uc.validateRegisterData(phone, password, name); err != nil {
		return nil, err
	}

	// Проверяем, существует ли пользователь
	existingUser, _ := uc.userRepo.GetUserByPhone(phone)
	if existingUser != nil {
		return nil, errors.New("user already exists")
	}

	// Хешируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Создаем пользователя
	user := &domain.User{
		Phone:    phone,
		Password: string(hashedPassword),
		Name:     name,
	}

	if err := uc.userRepo.CreateUser(user); err != nil {
		return nil, err
	}

	// Генерируем JWT токен
	token, err := uc.jwtHelper.GenerateToken(user.ID, user.Phone)
	if err != nil {
		return nil, err
	}

	// Очищаем пароль перед возвратом
	user.Password = ""

	return &AuthResponse{
		User:  user,
		Token: token,
	}, nil
}

// AuthenticateUser аутентифицирует пользователя
func (uc *AuthUseCase) AuthenticateUser(phone, password string) (*AuthResponse, error) {
	// Валидация входных данных
	if err := uc.validateLoginData(phone, password); err != nil {
		return nil, err
	}

	// Получаем пользователя
	user, err := uc.userRepo.GetUserByPhone(phone)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Проверяем пароль
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Генерируем JWT токен
	token, err := uc.jwtHelper.GenerateToken(user.ID, user.Phone)
	if err != nil {
		return nil, err
	}

	// Очищаем пароль перед возвратом
	user.Password = ""

	return &AuthResponse{
		User:  user,
		Token: token,
	}, nil
}

// GetUserByID получает пользователя по ID (для использования в middleware)
func (uc *AuthUseCase) GetUserByID(userID uint) (*domain.User, error) {
	user, err := uc.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}
	
	// Очищаем пароль
	user.Password = ""
	return user, nil
}

// validateRegisterData валидирует данные для регистрации
func (uc *AuthUseCase) validateRegisterData(phone, password, name string) error {
	if strings.TrimSpace(phone) == "" {
		return errors.New("phone is required")
	}
	
	if strings.TrimSpace(password) == "" {
		return errors.New("password is required")
	}
	
	if strings.TrimSpace(name) == "" {
		return errors.New("name is required")
	}

	// Валидация формата телефона (базовая)
	phoneRegex := regexp.MustCompile(`^\+?[1-9]\d{1,14}$`)
	if !phoneRegex.MatchString(phone) {
		return errors.New("invalid phone format")
	}

	// Валидация пароля (минимум 6 символов)
	if len(password) < 6 {
		return errors.New("password must be at least 6 characters long")
	}

	return nil
}

// validateLoginData валидирует данные для входа
func (uc *AuthUseCase) validateLoginData(phone, password string) error {
	if strings.TrimSpace(phone) == "" {
		return errors.New("phone is required")
	}
	
	if strings.TrimSpace(password) == "" {
		return errors.New("password is required")
	}

	return nil
}
