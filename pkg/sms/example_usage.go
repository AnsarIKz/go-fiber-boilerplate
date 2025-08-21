package sms

import (
	"fmt"
	"log"
)

// ExampleUsage демонстрирует использование TwilioSMSService
// Этот файл предназначен только для демонстрации и не должен использоваться в production
func ExampleUsage() {
	// Пример 1: Создание сервиса с явной конфигурацией
	config := &TwilioConfig{
		AccountSID: "your_account_sid",
		AuthToken:  "your_auth_token",
		From:       "+1234567890", // Ваш Twilio номер
	}

	smsService := NewTwilioSMSService(config)

	// Пример 2: Создание сервиса из переменных окружения
	// smsService := NewTwilioSMSServiceFromEnv()

	// Создание SMS сообщения
	sms := &SMSMessage{
		To:   "+0987654321", // Номер получателя в формате E.164
		Body: "Hello from Twilio SMS Service!",
	}

	// Отправка SMS
	err := smsService.SendSMS(sms)
	if err != nil {
		log.Printf("Failed to send SMS: %v", err)
		return
	}

	fmt.Println("SMS sent successfully!")
}

// CompareWithMailer показывает схожесть интерфейса с mailer
func CompareWithMailer() {
	// Twilio SMS Service
	smsService := NewTwilioSMSServiceFromEnv()
	sms := &SMSMessage{
		To:   "+1234567890",
		Body: "Your verification code is 123456",
	}
	// err := smsService.SendSMS(sms)

	// Mailer (для сравнения)
	// mailer := mailer.NewSMTPMailerFromEnv()
	// email := &mailer.EmailMessage{
	//     To:      []string{"user@example.com"},
	//     Subject: "Verification Code",
	//     Body:    "Your verification code is 123456",
	// }
	// err := mailer.SendEmail(email)

	// Оба сервиса имеют схожий интерфейс и могут быть легко заменены
	fmt.Printf("SMS Service: %T\n", smsService)
	fmt.Printf("SMS Message: %T\n", sms)
}
