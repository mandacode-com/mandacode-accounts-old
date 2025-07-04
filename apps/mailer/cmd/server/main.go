package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-playground/validator/v10"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
	kafkaserver "mandacode.com/accounts/mailer/cmd/server/kafka"
	mailapp "mandacode.com/accounts/mailer/internal/app/mail"
	"mandacode.com/accounts/mailer/internal/config"
	mailhandler "mandacode.com/accounts/mailer/internal/handler/mail"
)

func main() {
	// Initialize logger
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("failed to create logger: %v", err)
	}

	validator := validator.New()

	cfg, err := config.LoadConfig(validator)
	if err != nil {
		logger.Fatal("failed to load configuration", zap.Error(err))
	}

	// Initialize MailApp
	mailApp, err := mailapp.NewMailApp(cfg.Mail.Host, cfg.Mail.Port, cfg.Mail.Username, cfg.Mail.Password, cfg.Mail.Sender, logger)
	if err != nil {
		logger.Fatal("failed to create MailApp", zap.Error(err))
	}

	// Create Kafka mailReader (consumer)
	mailReader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:9094"},
		Topic:   "accounts.user.email_verification_request",
		GroupID: "mailer-service-group",
	})
	defer mailReader.Close()

	// Create MailHandler
	mailHandler := mailhandler.NewMailHandler(mailApp, validator)

	// Create Kafka server with reader and handler
	kafkaServer := kafkaserver.NewKafkaServer(
		logger,
		[]*kafkaserver.ReaderHandler{
			{
				Reader:  mailReader,
				Handler: mailHandler,
			},
		})

	// Start Kafka server
	if err := kafkaServer.Start(); err != nil {
		logger.Fatal("failed to start Kafka server", zap.Error(err))
	}

	// Graceful shutdown
	waitForShutdown(kafkaServer.Stop)
}

// Handle graceful shutdown
func waitForShutdown(stopFunc func() error) {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	<-signalChan
	log.Println("Received shutdown signal, stopping server...")
	if err := stopFunc(); err != nil {
		log.Printf("Error stopping server: %v\n", err)
	} else {
		log.Println("Server stopped gracefully")
	}
}
