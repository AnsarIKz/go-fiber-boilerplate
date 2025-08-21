package sms

import (
	"errors"
	"fmt"
	"nodabackend/pkg/env"

	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

// SMSMessage представляет SMS сообщение
type SMSMessage struct {
	To   string // Номер телефона получателя в формате E.164
	Body string // Текст сообщения
}

// SMSService интерфейс для отправки SMS
type SMSService interface {
	SendSMS(msg *SMSMessage) error
}

// TwilioSMSService реализация SMSService через Twilio API
type TwilioSMSService struct {
	client *twilio.RestClient
	from   string // Twilio номер телефона
}

// TwilioConfig конфигурация Twilio
type TwilioConfig struct {
	AccountSID string
	AuthToken  string
	From       string // Twilio номер телефона
}

// NewTwilioSMSService создает новый Twilio SMS сервис
func NewTwilioSMSService(config *TwilioConfig) *TwilioSMSService {
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: config.AccountSID,
		Password: config.AuthToken,
	})

	return &TwilioSMSService{
		client: client,
		from:   config.From,
	}
}

// NewTwilioSMSServiceFromEnv создает Twilio SMS сервис из переменных окружения
func NewTwilioSMSServiceFromEnv() *TwilioSMSService {
	accountSID := env.GetEnvOrDefault("TWILIO_ACCOUNT_SID", "")
	authToken := env.GetEnvOrDefault("TWILIO_AUTH_TOKEN", "")
	from := env.GetEnvOrDefault("TWILIO_PHONE_NUMBER", "")

	return NewTwilioSMSService(&TwilioConfig{
		AccountSID: accountSID,
		AuthToken:  authToken,
		From:       from,
	})
}

// SendSMS отправляет SMS через Twilio API
func (t *TwilioSMSService) SendSMS(msg *SMSMessage) error {
	if msg.To == "" {
		return errors.New("recipient phone number is required")
	}

	if msg.Body == "" {
		return errors.New("message body is required")
	}

	if t.from == "" {
		return errors.New("Twilio phone number is required")
	}

	// Валидация номера телефона (базовая проверка формата E.164)
	if !isValidPhoneNumber(msg.To) {
		return errors.New("invalid phone number format, expected E.164 format (e.g., +1234567890)")
	}

	// Параметры для отправки SMS
	params := &openapi.CreateMessageParams{}
	params.SetTo(msg.To)
	params.SetFrom(t.from)
	params.SetBody(msg.Body)

	// Отправляем SMS через Twilio API
	resp, err := t.client.Api.CreateMessage(params)
	if err != nil {
		return fmt.Errorf("failed to send SMS via Twilio: %w", err)
	}

	// Проверяем статус ответа
	if resp.Status == nil || *resp.Status != "queued" && *resp.Status != "sent" && *resp.Status != "delivered" {
		return fmt.Errorf("SMS failed to send, status: %s", getStatusString(resp.Status))
	}

	fmt.Printf("[TWILIO SMS] SMS sent successfully to: %s, SID: %s, Status: %s\n",
		msg.To, getSidString(resp.Sid), getStatusString(resp.Status))

	return nil
}

// getSidString safely извлекает строку из указателя на строку
func getSidString(sid *string) string {
	if sid == nil {
		return "unknown"
	}
	return *sid
}

// getStatusString safely извлекает строку из указателя на строку
func getStatusString(status *string) string {
	if status == nil {
		return "unknown"
	}
	return *status
}

// isValidPhoneNumber проверяет формат номера телефона (E.164)
func isValidPhoneNumber(phoneNumber string) bool {
	if len(phoneNumber) < 10 || len(phoneNumber) > 15 {
		return false
	}

	// Проверяем, что начинается с +
	if phoneNumber[0] != '+' {
		return false
	}

	// Проверяем, что остальная часть состоит только из цифр
	for i := 1; i < len(phoneNumber); i++ {
		if phoneNumber[i] < '0' || phoneNumber[i] > '9' {
			return false
		}
	}

	return true
}
