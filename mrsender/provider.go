package mrsender

import (
	"context"
	"net/textproto"
)

const (
	ContentTypePlain = "text/plain" // ContentTypePlain - простой текст
)

type (
	// MailProvider - провайдер для отправки электронных писем.
	MailProvider interface {
		SendMail(ctx context.Context, from string, to []string, header textproto.MIMEHeader, body string) error
	}

	// MessageProvider - провайдер для отправки сообщений в конкретный мессенджер.
	MessageProvider interface {
		SendToChat(ctx context.Context, chatKey, message string) error
	}
)
