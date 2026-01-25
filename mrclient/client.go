package mrclient

import (
	"context"
	"net/textproto"
)

type (
	// MailSender - отправитель электронных писем.
	MailSender interface {
		SendMail(ctx context.Context, from string, to []string, header textproto.MIMEHeader, body string) error
	}

	// MessengerSender - отправитель сообщений в конкретный мессенджер.
	MessengerSender interface {
		SendToChat(ctx context.Context, chatKey, message string) error
	}
)
