package mailer

import (
	"errors"
	"fmt"
	"nodabackend/pkg/env"
	"strings"

	"net/smtp"
)

// EmailMessage представляет email сообщение
type EmailMessage struct {
	To      []string
	Subject string
	Body    string
	IsHTML  bool
}

// Mailer интерфейс для отправки email
type Mailer interface {
	SendEmail(msg *EmailMessage) error
}

// SMTPMailer реализация Mailer через SMTP
type SMTPMailer struct {
	host     string
	port     int
	username string
	password string
	from     string
}

// SMTPConfig конфигурация SMTP
type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}

// NewSMTPMailer создает новый SMTP mailer
func NewSMTPMailer(config *SMTPConfig) *SMTPMailer {
	return &SMTPMailer{
		host:     config.Host,
		port:     config.Port,
		username: config.Username,
		password: config.Password,
		from:     config.From,
	}
}

// NewSMTPMailerFromEnv создает SMTP mailer из переменных окружения
func NewSMTPMailerFromEnv() *SMTPMailer {
	host := env.GetEnvOrDefault("SMTP_HOST", "smtp.gmail.com")
	port := env.GetEnvIntOrDefault("SMTP_PORT", 587)
	username := env.GetEnvOrDefault("SMTP_USERNAME", "")
	password := env.GetEnvOrDefault("SMTP_PASSWORD", "")
	from := env.GetEnvOrDefault("SMTP_FROM", username)

	return NewSMTPMailer(&SMTPConfig{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
		From:     from,
	})
}

// SendEmail отправляет email через SMTP
func (m *SMTPMailer) SendEmail(msg *EmailMessage) error {
	if len(msg.To) == 0 {
		return errors.New("recipients list is empty")
	}

	if msg.Subject == "" {
		return errors.New("subject is required")
	}

	if msg.Body == "" {
		return errors.New("body is required")
	}

	// Создаем SMTP клиент
	client, err := smtp.Dial(fmt.Sprintf("%s:%d", m.host, m.port))
	if err != nil {
		return fmt.Errorf("failed to connect to SMTP server: %w", err)
	}
	defer client.Close()

	// Если требуется аутентификация
	if m.username != "" && m.password != "" {
		auth := smtp.PlainAuth("", m.username, m.password, m.host)
		if err = client.Auth(auth); err != nil {
			return fmt.Errorf("SMTP authentication failed: %w", err)
		}
	}

	// Устанавливаем отправителя
	if err = client.Mail(m.from); err != nil {
		return fmt.Errorf("failed to set sender: %w", err)
	}

	// Устанавливаем получателей
	for _, to := range msg.To {
		if err = client.Rcpt(to); err != nil {
			return fmt.Errorf("failed to set recipient %s: %w", to, err)
		}
	}

	// Получаем writer для тела письма
	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("failed to get data writer: %w", err)
	}

	// Формируем заголовки
	contentType := "text/plain"
	if msg.IsHTML {
		contentType = "text/html"
	}

	// Пишем заголовки и тело
	_, err = fmt.Fprintf(w, "From: %s\r\n", m.from)
	if err != nil {
		return fmt.Errorf("failed to write From header: %w", err)
	}

	_, err = fmt.Fprintf(w, "To: %s\r\n", strings.Join(msg.To, ", "))
	if err != nil {
		return fmt.Errorf("failed to write To header: %w", err)
	}

	_, err = fmt.Fprintf(w, "Subject: %s\r\n", msg.Subject)
	if err != nil {
		return fmt.Errorf("failed to write Subject header: %w", err)
	}

	_, err = fmt.Fprintf(w, "Content-Type: %s; charset=UTF-8\r\n", contentType)
	if err != nil {
		return fmt.Errorf("failed to write Content-Type header: %w", err)
	}

	_, err = fmt.Fprintf(w, "\r\n%s", msg.Body)
	if err != nil {
		return fmt.Errorf("failed to write body: %w", err)
	}

	// Закрываем writer
	err = w.Close()
	if err != nil {
		return fmt.Errorf("failed to close data writer: %w", err)
	}

	// Отправляем QUIT
	err = client.Quit()
	if err != nil {
		return fmt.Errorf("failed to quit SMTP session: %w", err)
	}

	return nil
}




