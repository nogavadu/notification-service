package email

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/nogavadu/notification-service/internal/service"
	"gopkg.in/gomail.v2"
	"log/slog"
)

type emailService struct {
	log *slog.Logger
}

func New(
	log *slog.Logger,
) service.EmailService {
	return &emailService{
		log: log,
	}
}

func (s *emailService) SendMsg(ctx context.Context, to []string, subject, text string) error {
	s.log.Info(fmt.Sprintf("Sending email to %v", to))

	m := gomail.NewMessage()

	m.SetHeader("From", "test_golang@mail.ru")
	m.SetHeader("To", "test_golang@mail.ru")
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", text)

	d := gomail.NewDialer("smtp.mail.ru", 587, "test_golang@mail.ru", "sp2c20cRJDLbwi0hUS6J")
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	conn, err := d.Dial()
	if err != nil {
		s.log.Error(fmt.Sprintf("Dial error: %s", err.Error()))
		return err
	}
	defer conn.Close()

	if err = gomail.Send(conn, m); err != nil {
		s.log.Error(fmt.Sprintf("Sending error: %s", err.Error()))
		return err
	}
	return nil
}
