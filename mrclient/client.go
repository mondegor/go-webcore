package mrclient

import (
	"context"
	"net/textproto"
)

type (
	// MailSender - предоставляет метод для отправки электронных писем.
	MailSender interface {
		SendMail(ctx context.Context, from string, to []string, header textproto.MIMEHeader, body string) error
	}

	// MessengerSender - предоставляет метод для отправки сообщений в чат/мессенджер.
	MessengerSender interface {
		SendToChat(ctx context.Context, chatKey, message string) error
	}
)
