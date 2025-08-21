# NODA Backend - SMS Authentication System

Система аутентификации с уникальной логикой:

- **Заявки вместо регистрации** - пользователи подают заявки на регистрацию
- **Админское одобрение** - заявки проверяются и одобряются/отклоняются админом
- **SMS-авторизация** - авторизация по номеру телефона через SMS-код без пароля

## Архитектура

```
internal/
├── applications/          # 💎 Заявки на регистрацию
│   ├── submit.go          # Логика подачи заявки
│   ├── approve.go         # Логика одобрения заявки (для админа)
│   ├── reject.go          # Логика отклонения заявки (для админа)
│   └── routes.go          # Роуты для управления заявками
│
├── auth/                  # 💎 Аутентификация
│   ├── request_code.go    # Шаг 1: Запрос СМС-кода по номеру телефона
│   ├── verify_code.go     # Шаг 2: Проверка кода и выдача JWT
│   ├── middleware.go      # Middleware для защиты роутов
│   └── routes.go
│
├── users/                 # 💎 Пользователи
│   ├── user.go            # Модель пользователя (Phone, Status)
│   └── repository.go      # Методы для работы с БД
│

pkg/
├── jwthelper/jwthelper.go # JWT токены
├── otp/                   # 🍀 One-Time Password
│   └── otp.go             # Генерация и хранение кодов в Redis
└── sms/                   # 🍀 SMS Service
    └── client.go          # Интерфейс и реализация для отправки СМС
```

## Настройка

### 1. Переменные окружения

Создайте файл `.env` в корне проекта:

```env
# Конфигурация базы данных PostgreSQL
DB_HOST=localhost
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=nodabackend
DB_PORT=5432

# JWT секретный ключ для подписи токенов
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production

# Конфигурация Redis для OTP кодов
REDIS_URL=localhost:6379

# Конфигурация SMS провайдера
SMS_PROVIDER=mock # mock, twilio
TWILIO_ACCOUNT_SID=your-twilio-account-sid
TWILIO_AUTH_TOKEN=your-twilio-auth-token
TWILIO_FROM_NUMBER=+1234567890

# Порт сервера
PORT=3000

# Администратор по умолчанию (создается автоматически при первом запуске)
ADMIN_PHONE=+1234567890
ADMIN_NAME=System Administrator
```

### 2. Зависимости

```bash
# Установка зависимостей
go mod tidy

# Убедитесь, что PostgreSQL и Redis запущены
# Для быстрого старта можно использовать Docker Compose:
docker-compose up -d
```

### 3. Запуск

```bash
go run cmd/api/main.go
```

## API Endpoints

### Публичные роуты

#### Подача заявки на регистрацию

```http
POST /api/v1/applications/submit
Content-Type: application/json

{
  "phone": "+1234567890",
  "name": "John Doe"
}
```

#### Запрос SMS-кода для входа

```http
POST /api/v1/auth/request-code
Content-Type: application/json

{
  "phone": "+1234567890"
}
```

#### Проверка SMS-кода и получение JWT

```http
POST /api/v1/auth/verify-code
Content-Type: application/json

{
  "phone": "+1234567890",
  "code": "123456"
}
```

### Защищенные роуты (требуют JWT)

#### Получение профиля пользователя

```http
GET /api/v1/protected/profile
Authorization: Bearer <jwt-token>
```

### Админские роуты (требуют JWT + admin права)

#### Получение заявок на рассмотрении

```http
GET /api/v1/applications/admin/pending?limit=20&offset=0
Authorization: Bearer <admin-jwt-token>
```

#### Одобрение заявки

```http
POST /api/v1/applications/admin/approve
Authorization: Bearer <admin-jwt-token>
Content-Type: application/json

{
  "user_id": 1
}
```

#### Отклонение заявки

```http
POST /api/v1/applications/admin/reject
Authorization: Bearer <admin-jwt-token>
Content-Type: application/json

{
  "user_id": 1,
  "reason": "Недостаточно данных"
}
```

## Статусы пользователей

- `pending` - Заявка подана, ожидает одобрения
- `approved` - Заявка одобрена, пользователь может авторизоваться
- `rejected` - Заявка отклонена
- `blocked` - Пользователь заблокирован

## Паттерны проектирования (Go Best Practices)

### 🏗️ Архитектурные паттерны

- **Clean Architecture**: Разделение на domain/usecase/repository слои
- **Dependency Injection**: Интерфейсы в usecase, реализации в repository
- **Factory Pattern**: Функции `New...()` для создания объектов через интерфейсы

### 🔧 Порождающие паттерны

- **Functional Options**: Гибкая конфигурация (`WithTimeout`, `WithMaxConn`)
- **Builder Pattern**: Для сложных объектов с множеством параметров

### 🏛️ Структурные паттерны

- **Decorator**: HTTP middleware (auth, logging, metrics)
- **Interface Segregation**: Мелкие интерфейсы вместо больших

### ⚡ Конкурентные паттерны

- **Worker Pool**: Ограничение горутин для фоновых задач
- **Fan-out/Fan-in**: Параллельная обработка и сбор результатов
- **Context**: Отмена операций и передача метаданных

### 📋 Дополнительно

- **Error Handling**: Явная обработка ошибок, sentinel errors
- **Graceful Shutdown**: Корректное завершение работы

## Документация

📖 **Подробная документация**: [docs/README.md](./docs/README.md)

### Конвенции проекта

- [📝 Правила именования](./docs/NAMING_CONVENTIONS.md)
- [🏗️ Структура проекта](./docs/PROJECT_STRUCTURE.md)
- [🚀 Коммиты](./docs/COMMIT_CONVENTIONS.md)

## Особенности реализации

1. **Безопасность OTP**: Коды действуют 5 минут, ограничение на количество попыток
2. **SMS уведомления**: Автоматические уведомления о статусе заявки
3. **Админская панель**: Управление заявками и пользователями
4. **JWT токены**: 24 часа жизни с возможностью обновления
5. **Автомиграции**: База данных настраивается автоматически при первом запуске

## Структура БД

### Таблица `users`

- `id` - Уникальный идентификатор
- `phone` - Номер телефона (уникальный)
- `name` - Имя пользователя
- `status` - Статус заявки/пользователя
- `is_admin` - Флаг администратора
- `created_at`, `updated_at`, `deleted_at` - Временные метки
