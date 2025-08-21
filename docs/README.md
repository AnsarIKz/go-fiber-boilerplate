# Go Project Documentation

## 📋 Документация проекта

### 🎯 Конвенции и стандарты

#### [📝 NAMING_CONVENTIONS.md](./NAMING_CONVENTIONS.md)

Правила именования для Go кода:

- Пакеты, структуры, интерфейсы
- Переменные, константы, функции
- Файлы и директории
- Go-идиомы и anti-patterns

#### [🏗️ PROJECT_STRUCTURE.md](./PROJECT_STRUCTURE.md)

Структура Go проекта:

- Clean Architecture слои
- Domain-Driven Design организация
- Примеры кода для каждого слоя
- Dependency injection паттерны

#### [🚀 COMMIT_CONVENTIONS.md](./COMMIT_CONVENTIONS.md)

Conventional Commits для Go:

- Формат коммитов
- Типы и scope'ы
- Примеры хороших коммитов
- Автоматизация с husky

### 🔧 Дополнительная документация

#### [API.md](./API.md) (предполагается создать)

- REST API спецификация
- Примеры запросов/ответов
- Аутентификация и авторизация

#### [DEPLOYMENT.md](./DEPLOYMENT.md) (предполагается создать)

- Инструкции по развертыванию
- Docker конфигурация
- CI/CD pipeline

### 📚 Паттерны проектирования Go

#### Архитектурные

- **Clean Architecture**: Разделение ответственности
- **DDD**: Domain-Driven Design
- **Hexagonal Architecture**: Порты и адаптеры

#### Go-специфичные

- **Functional Options**: Гибкая конфигурация
- **Interface Composition**: Композиция интерфейсов
- **Error Handling**: Явное возвращение ошибок
- **Context**: Отмена и таймауты

### 🛠️ Инструменты

#### Качество кода

- `golangci-lint` - линтер
- `go fmt` - форматирование
- `go vet` - статический анализ

#### Тестирование

- `go test` - юнит и интеграционные тесты
- `testify` - assertions и mocks
- `httptest` - тестирование HTTP

#### Документация

- `swag` - Swagger генерация
- `godoc` - документация кода

## 🚀 Быстрый старт

1. **Прочитайте** [NAMING_CONVENTIONS.md](./NAMING_CONVENTIONS.md)
2. **Изучите** [PROJECT_STRUCTURE.md](./PROJECT_STRUCTURE.md)
3. **Следуйте** [COMMIT_CONVENTIONS.md](./COMMIT_CONVENTIONS.md)
4. **Соблюдайте** паттерны из основного README

## 📖 Полезные ссылки

- [Effective Go](https://golang.org/doc/effective_go.html)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Conventional Commits](https://conventionalcommits.org/)
- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
