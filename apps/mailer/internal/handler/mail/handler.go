package mailhandler

import (
	"github.com/go-playground/validator/v10"
	kafka "github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"
	kafkaserver "mandacode.com/accounts/mailer/cmd/server/kafka"
	mailapp "mandacode.com/accounts/mailer/internal/app/mail"
	mailerv1 "mandacode.com/accounts/proto/mailer/v1"
)

type MailHandler struct {
	MailApp   mailapp.MailApp
	validator *validator.Validate
}

// HandleMessage implements kafkaserver.KafkaHandler.
func (h *MailHandler) HandleMessage(m kafka.Message) error {
	event := &mailerv1.EmailVerificationEvent{}
	if err := proto.Unmarshal(m.Value, event); err != nil {
		return err
	}
	if err := h.MailApp.SendEmailVerificationMail(event.Email, event.VerificationLink); err != nil {
		return err
	}

	return nil
}

func NewMailHandler(mailApp mailapp.MailApp, validator *validator.Validate) kafkaserver.KafkaHandler {
	return &MailHandler{
		MailApp:   mailApp,
		validator: validator,
	}
}
