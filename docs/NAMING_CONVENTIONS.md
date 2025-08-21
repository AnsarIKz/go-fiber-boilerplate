# Go Naming Conventions

## 📝 Общие правила

### Go Style Guide

- **Пакеты**: строчные буквы, без подчеркиваний (`user`, `auth`, `repository`)
- **Интерфейсы**: существительные или прилагательные (`Repository`, `Service`, `Handler`)
- **Структуры**: существительные с заглавной буквы (`User`, `AuthService`)
- **Методы**: глаголы с заглавной буквы (`GetUser`, `CreateToken`, `ValidateCode`)
- **Функции**: глаголы с заглавной буквы (`NewUser`, `ParseToken`)
- **Переменные**: camelCase (`userID`, `phoneNumber`, `isActive`)
- **Константы**: SCREAMING_SNAKE_CASE (`MAX_RETRIES`, `JWT_SECRET`)

## 🏗️ Архитектурные компоненты

### Файлы и директории

```
internal/
├── users/                    # строчные буквы, множественное число
│   ├── domain/
│   │   ├── user.go          # существительное
│   │   └── interfaces.go    # множественное число для коллекций
│   ├── usecase/
│   │   └── create_user.go   # snake_case для файлов с несколькими словами
│   ├── repository/
│   │   └── postgres_user_repository.go
│   └── delivery/
│       └── http/
│           └── user_handler.go

pkg/
├── database/                # строчные буквы
│   └── postgres.go
└── jwthelper/               # строчные буквы
    └── jwthelper.go         # повторение имени пакета
```

### Интерфейсы

```go
// Хорошие имена
type UserRepository interface {
    GetByID(id int) (*User, error)
    GetByPhone(phone string) (*User, error)
    Create(user *User) error
    Update(user *User) error
    Delete(id int) error
}

type AuthService interface {
    GenerateToken(userID int) (string, error)
    ValidateToken(token string) (int, error)
}

// Избегайте
type IRepository interface {} // Нет префикса I
type UserInterface interface {} // Слишком общее
```

### Структуры и методы

```go
// Конструкторы
func NewUserRepository(db *sql.DB) *UserRepository { ... }
func NewAuthService(repo UserRepository) *AuthService { ... }

// Структуры
type User struct {
    ID        int       `json:"id"`
    Phone     string    `json:"phone"`
    Name      string    `json:"name"`
    Status    string    `json:"status"`
    IsAdmin   bool      `json:"is_admin"`
    CreatedAt time.Time `json:"created_at"`
}

// Методы
func (u *User) Validate() error {
    return validationErrors{...}
}

func (u *User) IsActive() bool {
    return u.Status == "approved"
}
```

## 🔧 Переменные и константы

### Переменные

```go
// Локальные
var (
    userID    int
    userName  string
    isActive  bool
)

// Глобальные (избегайте, используйте dependency injection)
var (
    db     *sql.DB
    logger *log.Logger
)
```

### Константы

```go
const (
    // Статусы пользователей
    StatusPending   = "pending"
    StatusApproved  = "approved"
    StatusRejected  = "rejected"
    StatusBlocked   = "blocked"

    // Ограничения
    MaxRetries     = 3
    TokenLifetime  = 24 * time.Hour
    CodeLength     = 6
)
```

## 🎯 Специальные случаи

### Context

```go
ctx := context.Background()
ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
defer cancel()
```

### Errors

```go
// Sentinel errors
var (
    ErrUserNotFound    = errors.New("user not found")
    ErrInvalidCode     = errors.New("invalid verification code")
    ErrTooManyAttempts = errors.New("too many attempts")
)

// Error wrapping
if err != nil {
    return fmt.Errorf("failed to create user: %w", err)
}
```

### HTTP Handlers

```go
// Имена функций
func HandleGetUser(c *gin.Context) { ... }
func HandleCreateUser(c *gin.Context) { ... }
func HandleUpdateUser(c *gin.Context) { ... }
func HandleDeleteUser(c *gin.Context) { ... }
```

## 🚫 Anti-patterns

### Избегайте

```go
// Нет
type IUserRepository interface {}  // I-префикс
var user_id int                    // snake_case
func getUser() {}                  // строчная буква
func new_user() {}                 // snake_case
type user struct {}                // строчная буква
const MAX_RETRIES = 3              // нет необходимости в package scope
```

### Предпочитайте

```go
// Да
type UserRepository interface {}
var userID int
func GetUser() {}
func NewUser() {}
type User struct {}
const MaxRetries = 3
```
