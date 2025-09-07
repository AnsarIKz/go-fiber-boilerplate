# Go Fiber Boilerplate

This is a boilerplate for Go applications built with the [Fiber](https://gofiber.io/) framework. It follows the principles of Clean Architecture and comes with a ready-to-use setup for PostgreSQL, Redis, and JWT-based authentication.

## 🚀 Tech Stack

- **Language:** [Go](https://golang.org/)
- **Web Framework:** [Fiber](https://gofiber.io/)
- **Database:** [PostgreSQL](https://www.postgresql.org/)
- **Cache:** [Redis](https://redis.io/)
- **ORM:** [GORM](https://gorm.io/)
- **Containerization:** [Docker](https://www.docker.com/)

## 🚧 Project To-Do List

- [x] **Authentication**: Basic JWT authentication (login/register).
- [x] **OTP Service**: Complete OTP functionality with Redis (generation, validation, expiry).
- [x] **Email Service**: SMTP mailer for sending emails.
- [x] **SMS Service**: Twilio integration for SMS sending.
- [x] **Redis Integration**: Ready-to-use Redis client with singleton pattern.
- [ ] **Refresh Tokens**: Implement refresh tokens for better security.
- [ ] **OTP Integration**: Connect OTP service to auth endpoints (login/register/recovery).
- [ ] **Extended Auth**: Password recovery, social login (OAuth2).
- [ ] **File Management**: Upload, storage, and processing.
  - [ ] File upload/storage system
  - [ ] Image processing and optimization
  - [ ] CDN integration for static assets
- [ ] **Database Optimization**: Performance and reliability.
  - [ ] Database indexing strategy
  - [ ] Database backup strategy
  - [ ] Caching strategy with Redis
- [ ] **DevOps & Infrastructure**: Production-ready deployment.
  - [ ] CI/CD Pipeline (GitHub Actions/GitLab CI)
  - [ ] Environment management (dev/staging/prod)
  - [ ] Load balancing configuration
- [ ] **Monitoring & Observability**: Metrics, logging, tracing, and health checks.
  - [ ] Structured logging (e.g., `logrus`, `zap`)
  - [ ] Metrics collection (Prometheus)
  - [ ] Error tracking and alerting
  - [ ] Request tracing and performance monitoring
  - [ ] Health check endpoints
- [ ] **Configuration**: Advanced configuration management (e.g., Viper).
- [ ] **Testing**: Add Unit and Integration tests.
- [ ] **API Documentation**: Auto-generated docs with Swagger/OpenAPI.

## ⚙️ Getting Started

### Prerequisites

- [Docker](https://www.docker.com/get-started)
- [Go](https://golang.org/dl/) (v1.18+)

### Steps

1.  **Clone the repository:**

    ```bash
    git clone <YOUR_REPOSITORY_URL>
    cd nodabackend
    ```

2.  **Create `.env` file:**
    Create a `.env` file in the root directory and add the following configuration.

    ```env
    # Database
    DB_HOST=localhost
    DB_PORT=5432
    DB_USER=postgres
    DB_PASSWORD=postgres
    DB_NAME=noda
    DB_SSLMODE=disable

    # Redis
    REDIS_HOST=localhost
    REDIS_PORT=6379
    REDIS_PASSWORD=

    # JWT
    JWT_SECRET=your_super_secret_key
    JWT_EXPIRATION_HOURS=72

    # SMTP (for email notifications)
    SMTP_HOST=smtp.example.com
    SMTP_PORT=587
    SMTP_USER=user@example.com
    SMTP_PASSWORD=password
    SMTP_SENDER=sender@example.com
    ```

3.  **Run services:**
    This will start the PostgreSQL and Redis containers.

    ```bash
    docker-compose up -d
    ```

4.  **Run the application:**
    ```bash
    go mod tidy
    go run cmd/api/main.go
    ```

The server will be running on `http://localhost:3000`.
