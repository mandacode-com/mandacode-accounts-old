package mailer

import (
	"context"

	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
	mailerv1 "mandacode.com/accounts/proto/mailer/v1"
	mailerdomain "mandacode.com/accounts/user/internal/domain/port/mailer"
)

type Mailer struct {
	writer *kafka.Writer
}

// SendEmailVerificationMail implements mailerdomain.Mailer.
func (m *Mailer) SendEmailVerificationMail(email string, verificationLink string) error {
	event := &mailerv1.EmailVerificationEvent{
		Email:            email,
		VerificationLink: verificationLink,
		EventTime:        timestamppb.Now(),
	}
	// Marshal the event to protobuf bytes
	data, err := proto.Marshal(event)
	if err != nil {
		return err
	}

	// Create a message to send to Kafka
	message := kafka.Message{
		Key:   []byte(email),
		Value: data,
	}

	return m.writer.WriteMessages(context.Background(), message)
}

// NewMailer creates a new Mailer instance with the provided Kafka writer.
func NewMailer(writer *kafka.Writer) mailerdomain.Mailer {
	return &Mailer{
		writer: writer,
	}
}
