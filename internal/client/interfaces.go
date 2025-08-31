package client

import "context"

type NotificationClientInterface interface {
	SendEmail(ctx context.Context, req SendEmailRequest) error
}
