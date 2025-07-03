package mailhandler

import (
	"encoding/json"

	kafka "github.com/segmentio/kafka-go"
	kafkaserver "mandacode.com/accounts/mailer/cmd/server/kafka"
	mailapp "mandacode.com/accounts/mailer/internal/app/mail"
	mailhandlerdto "mandacode.com/accounts/mailer/internal/handler/mail/dto"
)

type MailHandler struct {
	// MailApp is the interface for sending emails.
	MailApp mailapp.MailApp
}

// HandleMessage implements kafkaserver.KafkaHandler.
func (h *MailHandler) HandleMessage(m kafka.Message) error {
	var dto mailhandlerdto.MailVerificationRequest
	if err := json.Unmarshal(m.Value, &dto); err != nil {
		return err
	}

	// Send email verification mail using the MailApp.
	if err := h.MailApp.SendEmailVerificationMail(dto.Email, dto.VerificationLink); err != nil {
		return err
	}

	return nil
}

func NewMailHandler(mailApp mailapp.MailApp) kafkaserver.KafkaHandler {
	return &MailHandler{
		MailApp: mailApp,
	}
}
