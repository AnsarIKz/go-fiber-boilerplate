# Go Project Structure Conventions

## 📁 Общая структура проекта

### Стандартная структура

```
.
├── cmd/                          # Точки входа приложения
│   └── api/                      # HTTP API сервер
│       └── main.go               # Главная функция
├── internal/                     # Приватный код приложения
│   ├── auth/                     # Домен: аутентификация
│   │   ├── domain/               # Бизнес-логика и модели
│   │   │   ├── model.go          # Структуры данных
│   │   │   └── interfaces.go     # Интерфейсы репозиториев
│   │   ├── usecase/              # Бизнес-сценарии (интеракторы)
│   │   │   └── auth_usecase.go   # Логика аутентификации
│   │   ├── repository/           # Реализация доступа к данным
│   │   │   └── postgres_auth_repository.go
│   │   └── delivery/             # HTTP handlers (контроллеры)
│   │       └── http/
│   │           └── auth_handler.go
│   ├── users/                    # Домен: пользователи
│   │   └── ...                   # Та же структура
│   ├── applications/             # Домен: заявки
│   │   └── ...                   # Та же структура
│   └── shared/                   # Общий код для всех доменов
│       └── middleware/           # HTTP middleware
├── pkg/                          # Публичные пакеты
│   ├── database/                 # Подключение к БД
│   │   └── postgres.go
│   ├── jwthelper/                # JWT утилиты
│   │   └── jwthelper.go
│   ├── otp/                      # One-time passwords
│   │   └── otp.go
│   ├── sms/                      # SMS сервис
│   │   └── client.go
│   └── logger/                   # Логирование
│       └── logger.go
├── configs/                      # Конфигурационные файлы
│   ├── config.go                 # Структуры конфигурации
│   └── local.yaml                # Локальная конфигурация
├── migrations/                   # SQL миграции
│   ├── 001_initial_schema.up.sql
│   └── 001_initial_schema.down.sql
├── docker-compose.yml            # Docker окружение
├── Dockerfile                    # Docker образ
├── go.mod                        # Go модуль
├── go.sum                        # Checksums зависимостей
├── README.md                     # Документация
├── .env.example                  # Пример переменных окружения
└── docs/                         # Дополнительная документация
    ├── API.md                    # API документация
    ├── NAMING_CONVENTIONS.md     # Правила именования
    └── PROJECT_STRUCTURE.md      # Эта структура
```

## 🏗️ Принципы организации

### 1. Domain-Driven Design (DDD)

```
internal/
├── auth/           # Ограниченный контекст (Bounded Context)
├── users/          # Ограниченный контекст
└── applications/   # Ограниченный контекст
```

### 2. Clean Architecture Layers

```
domain/       # Entities (бизнес-сущности)
usecase/      # Use Cases (бизнес-сценарии)
repository/   # Interface Adapters (адаптеры)
delivery/     # Interface Adapters (HTTP handlers)
```

### 3. Dependency Direction

```
delivery → usecase → domain ← repository
     ↓                           ↓
  HTTP API               PostgreSQL/Redis
```

### 4. Прагматичный подход к Clean Architecture

Несмотря на все преимущества, Clean Architecture не является серебряной пулей. Применяйте ее осознанно:
Главный принцип — **не создавайте лишних слоев абстракции**, если они не решают конкретную проблему сложности.

## 📦 Слои и их ответственность

### Domain Layer (Бизнес-логика)

- **Структуры**: `User`, `Application`, `AuthToken`
- **Интерфейсы**: `UserRepository`, `AuthService`
- **Ошибки**: `ErrUserNotFound`, `ErrInvalidToken`
- **Валидация**: Бизнес-правила

### Use Case Layer (Сценарии использования)

- **Структуры**: `CreateUserUseCase`, `AuthUseCase`
- **Логика**: Оркестрация вызовов к repository
- **DTO**: Входные/выходные структуры для API

### Repository Layer (Доступ к данным)

- **Реализации**: `PostgresUserRepository`, `RedisOTPRepository`
- **Интерфейсы**: Реализуют интерфейсы из domain
- **Запросы**: SQL, Redis команды

### Delivery Layer (HTTP API)

- **Handlers**: `UserHandler`, `AuthHandler`
- **Middleware**: `AuthMiddleware`, `LoggingMiddleware`
- **Routing**: Настройка маршрутов

## 🔧 Правила для каждого слоя

### Domain

```go
// domain/user.go
package domain

type User struct {
    ID        int
    Phone     string
    Name      string
    Status    string
    CreatedAt time.Time
}

type UserRepository interface {
    GetByID(id int) (*User, error)
    GetByPhone(phone string) (*User, error)
    Create(user *User) error
    Update(user *User) error
}
```

### Use Case

```go
// usecase/user_usecase.go
package usecase

type UserUseCase struct {
    userRepo domain.UserRepository
}

func NewUserUseCase(repo domain.UserRepository) *UserUseCase {
    return &UserUseCase{userRepo: repo}
}

func (uc *UserUseCase) CreateUser(phone, name string) error {
    // Бизнес-логика здесь
    user := &domain.User{Phone: phone, Name: name}
    return uc.userRepo.Create(user)
}
```

### Repository

```go
// repository/postgres_user_repository.go
package repository

type PostgresUserRepository struct {
    db *sql.DB
}

func NewPostgresUserRepository(db *sql.DB) *PostgresUserRepository {
    return &PostgresUserRepository{db: db}
}

func (r *PostgresUserRepository) Create(user *domain.User) error {
    // SQL запросы здесь
    return nil
}
```

### Delivery (HTTP Handler)

```go
// delivery/http/user_handler.go
package http

type UserHandler struct {
    userUseCase *usecase.UserUseCase
}

func NewUserHandler(uc *usecase.UserUseCase) *UserHandler {
    return &UserHandler{userUseCase: uc}
}

func (h *UserHandler) CreateUser(c *gin.Context) {
    // HTTP логика здесь
    // Парсинг JSON, вызов usecase, возврат ответа
}
```
