package service

import "context"

type EmailService interface {
	SendMsg(ctx context.Context, to []string, subject, text string) error
}
