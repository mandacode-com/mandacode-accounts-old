package userevent

import (
	"context"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
	usereventv1 "mandacode.com/accounts/proto/user/event/v1"
	usereventdomain "mandacode.com/accounts/user/internal/domain/port/user_event"
)

type UserEventService struct {
	usercreationFailed *kafka.Writer
	archiveUser        *kafka.Writer
	deleteUser         *kafka.Writer
}

// UserCreationFailed implements usereventdomain.UserEventService.
func (u *UserEventService) UserCreationFailed(userID uuid.UUID) error {
	event := &usereventv1.UserCreationFailedEvent{
		UserId:    userID.String(),
		EventTime: timestamppb.Now(),
	}
	// Marshal the event to protobuf bytes
	data, err := proto.Marshal(event)
	if err != nil {
		return err
	}

	// Create a message to send to Kafka
	message := kafka.Message{
		Key:   []byte(userID.String()),
		Value: data,
	}

	return u.usercreationFailed.WriteMessages(context.Background(), message)
}

// ArchiveUser implements usereventdomain.UserEventService.
func (u *UserEventService) ArchiveUser(userID uuid.UUID) error {
	event := &usereventv1.UserArchivedEvent{
		UserId:    userID.String(),
		EventTime: timestamppb.Now(),
	}
	// Marshal the event to protobuf bytes
	data, err := proto.Marshal(event)
	if err != nil {
		return err
	}

	// Create a message to send to Kafka
	message := kafka.Message{
		Key:   []byte(userID.String()),
		Value: data,
	}

	return u.archiveUser.WriteMessages(context.Background(), message)
}

// DeleteUser implements usereventdomain.UserEventService.
func (u *UserEventService) DeleteUser(userID uuid.UUID) error {
	event := &usereventv1.UserDeletedEvent{
		UserId:    userID.String(),
		EventTime: timestamppb.Now(),
	}
	// Marshal the event to protobuf bytes
	data, err := proto.Marshal(event)
	if err != nil {
		return err
	}

	// Create a message to send to Kafka
	message := kafka.Message{
		Key:   []byte(userID.String()),
		Value: data,
	}

	return u.deleteUser.WriteMessages(context.Background(), message)
}

// NewMailer creates a new Mailer instance with the provided Kafka writer.
func NewUserEventService(
	userCreationFailed *kafka.Writer,
	archiveUser *kafka.Writer,
	deleteUser *kafka.Writer,
) usereventdomain.UserEventService {
	return &UserEventService{
		usercreationFailed: userCreationFailed,
		archiveUser:        archiveUser,
		deleteUser:         deleteUser,
	}
}
