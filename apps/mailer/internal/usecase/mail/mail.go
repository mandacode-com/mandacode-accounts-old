package mail

import (
	"bytes"
	"html/template"
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"gopkg.in/gomail.v2"
	maildomain "mandacode.com/accounts/mailer/internal/domain/usecase/mail"
)

type mailUsecase struct {
	dialer              *gomail.Dialer
	verifyEmailTemplate *template.Template
	logger              *zap.Logger
	username            string
	sender              string
}

// SendEmailVerificationMail implements MailApp.
func (m *mailUsecase) SendEmailVerificationMail(email string, link string) error {
	data := struct {
		Link string
	}{
		Link: link,
	}

	var body bytes.Buffer
	if err := m.verifyEmailTemplate.Execute(&body, data); err != nil {
		m.logger.Error("failed to execute email template", zap.Error(err), zap.String("to", email))
		return err
	}

	msg := gomail.NewMessage()
	msg.SetHeader("From", msg.FormatAddress(m.username, m.sender))
	msg.SetHeader("To", email)
	msg.SetHeader("Subject", "[Mandacode] Email Verification")
	msg.SetBody("text/html", body.String())

	if err := m.dialer.DialAndSend(msg); err != nil {
		m.logger.Error("failed to send email", zap.Error(err), zap.String("to", email))
		return err
	}

	m.logger.Info("email sent successfully", zap.String("to", email))
	return nil
}

// NewMailApp creates a new instance of MailApp with the provided SMTP configuration.
func NewMailApp(host string, port int, username, password, sender string, logger *zap.Logger) (maildomain.MailUsecase, error) {
	dialer := gomail.NewDialer(host, port, username, password)
	cwd, err := os.Getwd()
	tmplPath := filepath.Join(cwd, "internal", "template", "verify_email.html")
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		logger.Error("failed to parse email template", zap.Error(err))
		return nil, err
	}

	return &mailUsecase{
		dialer:              dialer,
		verifyEmailTemplate: tmpl,
		logger:              logger,
		username:            username,
		sender:              sender,
	}, nil
}
