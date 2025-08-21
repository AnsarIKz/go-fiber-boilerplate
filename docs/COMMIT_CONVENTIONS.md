# Go Project Commit Conventions

## 📝 Conventional Commits для Go проектов

### Формат коммита

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

### Типы коммитов

#### 🚀 Основные типы

- **`feat:`** - Новая функциональность
- **`fix:`** - Исправление бага
- **`docs:`** - Изменения в документации
- **`style:`** - Форматирование, точки с запятой, отступы (без изменения логики)
- **`refactor:`** - Изменение кода без исправления багов или добавления фич
- **`test:`** - Добавление или обновление тестов
- **`chore:`** - Обслуживающие задачи (обновление зависимостей, конфигурация)

#### 🔧 Дополнительные типы для Go

- **`perf:`** - Улучшение производительности
- **`ci:`** - Изменения в CI/CD конфигурации
- **`build:`** - Изменения, влияющие на систему сборки или внешние зависимости
- **`revert:`** - Откат предыдущего коммита

### Примеры коммитов

#### ✅ Хорошие примеры

```bash
# Новая функциональность
feat: add SMS verification code endpoint
feat(auth): implement JWT token validation middleware
feat(users): add user registration application system

# Исправления
fix: handle empty phone number in SMS sending
fix(auth): correct JWT token expiration time validation
fix(repository): fix PostgreSQL connection pool exhaustion

# Документация
docs: update API documentation for auth endpoints
docs: add README section about project structure

# Рефакторинг
refactor: extract user validation logic to separate function
refactor(auth): simplify JWT token generation
refactor: rename UserRepository to UserRepo for consistency

# Тесты
test: add unit tests for user usecase
test(auth): add integration tests for JWT middleware
test: add benchmark tests for database queries

# Обслуживание
chore: update Go modules to latest versions
chore: add Docker configuration for development
chore: configure golangci-lint with custom rules
```

#### ❌ Плохие примеры

```bash
# Слишком общие
fix bug
update code
changes

# На русском
feat: добавил новую функцию
fix: исправил ошибку

# Слишком подробные
feat: add a new SMS verification code endpoint that handles the logic for sending SMS codes to users when they request authentication via their phone number and store the codes in Redis with a 5-minute expiration time

# Неправильный тип
feat: fix typo in README
test: update dependencies
```

## 🎯 Scope (Область применения)

### Для Go проектов

- **`api`**: HTTP API endpoints и handlers
- **`auth`**: Аутентификация и авторизация
- **`db`**: База данных и репозитории
- **`config`**: Конфигурация приложения
- **`docker`**: Docker файлы и конфигурация
- **`deps`**: Зависимости и модули

### Примеры с scope

```bash
feat(api): add user profile endpoint
fix(auth): handle expired JWT tokens gracefully
refactor(db): optimize user query performance
test(db): add integration tests for user repository
docs(api): update Swagger documentation
```

## 📋 Тело коммита (Body)

### Когда использовать

- **Breaking changes**: Описывать изменения, ломающие обратную совместимость
- **Сложные изменения**: Объяснять почему и как было сделано изменение
- **Связанные задачи**: Указывать ссылки на issues, PR

### Формат

```
feat: add user registration endpoint

- Add POST /api/v1/applications/submit endpoint
- Validate phone number format and name length
- Store application in pending status
- Send SMS notification to admin

BREAKING CHANGE: API now requires phone verification
```

## 🔗 Footer (Подвал)

### Стандартные форматы

- **BREAKING CHANGE:** - для изменений, ломающих совместимость
- **Closes #123** - для закрытия issue
- **Related to #456** - для связанных задач

### Примеры

```bash
feat: implement user application system

Closes #42
BREAKING CHANGE: User registration now requires admin approval
```

```bash
fix: handle database connection timeout

Related to #123
```

## 🚀 Полезные команды

### Проверка формата коммитов

```bash
# Установка commitlint
npm install -g @commitlint/cli @commitlint/config-conventional

# Создание конфигурации .commitlintrc.json
{
  "extends": ["@commitlint/config-conventional"]
}
```

### Автоматическое форматирование

```bash
# Установка commitizen
npm install -g commitizen cz-conventional-changelog

# Использование
git cz
```

### Хуки pre-commit

```bash
# Установка husky и lint-staged
npx husky-init && npm install
npx husky add .husky/pre-commit "golangci-lint run"
npx husky add .husky/commit-msg "npx commitlint --edit \$1"
```

## 📊 Workflow пример

```bash
# Начинаем новую фичу
git checkout -b feat/user-registration

# Делаем изменения
# ... кодинг ...

# Коммитим по частям
git add internal/users/
git commit -m "feat(users): add user domain models and interfaces"

git add internal/applications/
git commit -m "feat(applications): implement application submission usecase"

git add internal/delivery/http/
git commit -m "feat(api): add application submission endpoint"

# Отправляем на review
git push origin feat/user-registration
```

## 🎯 Go-специфичные рекомендации

### Структура коммитов для Go проектов

1. **API changes**: `feat(api): ...` или `fix(api): ...`
2. **Database**: `feat(db): ...` или `refactor(db): ...`
3. **Dependencies**: `chore(deps): ...`
4. **Configuration**: `chore(config): ...`
5. **Tests**: `test: ...` с указанием области

### Модульные коммиты

- **Маленькие**: Каждый коммит решает одну задачу
- **Атомарные**: Коммит можно откатить без проблем
- **Описательные**: Название объясняет что именно сделано
