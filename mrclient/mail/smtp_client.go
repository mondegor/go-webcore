package mail

import (
	"bytes"
	"context"
	"net/smtp"
	"net/textproto"
	"strings"

	"github.com/mondegor/go-core/errors"
	"github.com/mondegor/go-core/mrtrace"
)

const (
	smtpClientName = "SmtpClient"
)

type (
	// SMTPClient - адаптер для отправки электронных писем через SMTP.
	SMTPClient struct {
		address string
		auth    smtp.Auth
		tracer  mrtrace.Tracer
	}
)

// NewSMTPClient - создаёт и настраивает SMTP-клиент для отправки электронных писем.
// Параметры:
//   - host - SMTP-сервер;
//   - port - порт сервера;
//   - username и password - учётные данные для аутентификации;
//
// Параметр tracer используется для трассировки операций отправки писем.
func NewSMTPClient(
	host, port, username, password string,
	tracer mrtrace.Tracer,
) *SMTPClient {
	return &SMTPClient{
		address: host + ":" + port,
		auth:    smtp.PlainAuth("", username, password, host),
		tracer:  tracer,
	}
}

// SendMail - отправляет электронное письмо через SMTP-сервер.
// Параметры:
//   - from - email отправителя;
//   - to - список email получателей;
//   - header - MIME-заголовки письма;
//   - body - тело письма.
//
// Автоматически добавляет заголовки From, To и Subject, если они отсутствуют в header.
// Возвращает ошибку ErrInternalIncorrectInputData, если from или to пусты.
// Возвращает ошибку ErrInternalServiceOperationFailed при сбое отправки письма.
func (c *SMTPClient) SendMail(ctx context.Context, from string, to []string, header textproto.MIMEHeader, body string) error {
	if from == "" {
		return errors.ErrInternalIncorrectInputData.WithDetails("from address is empty")
	}

	if len(to) == 0 || to[0] == "" {
		return errors.ErrInternalIncorrectInputData.WithDetails("to address is empty")
	}

	var buf bytes.Buffer

	// если в заголовке отсутствует адрес отправителя
	if len(header) == 0 || header.Get("From") == "" {
		buf.WriteString("From: ")
		buf.WriteString(from)
		buf.WriteString("\r\n")
	}

	// если в заголовке отсутствует адрес получателя
	if len(header) == 0 || header.Get("To") == "" {
		buf.WriteString("To: ")
		buf.WriteString(to[0])
		buf.WriteString("\r\n")
	}

	// если в заголовке отсутствует тема письма
	if len(header) == 0 || header.Get("Subject") == "" {
		buf.WriteString("Subject: ")
		buf.WriteString(defaultMessageSubject)
		buf.WriteString("\r\n")
	}

	if len(header) > 0 {
		for canonicalKey := range header {
			headerValue := header[canonicalKey]

			// если значение у какого либо заголовка окажется пустым, то этот заголовок не будет сформирован
			if len(headerValue) == 0 || headerValue[0] == "" {
				continue
			}

			buf.WriteString(canonicalKey)
			buf.WriteString(": ")
			buf.WriteString(headerValue[0])
			buf.WriteString("\r\n")
		}
	}

	buf.WriteString("\r\n")
	buf.WriteString(body)

	c.tracer.Trace(
		ctx,
		"source", smtpClientName,
		"cmd", "SendMail",
		"SMTP-Header", buf.String(),
		"SMTP-From", from,
		"SMTP-To", strings.Join(to, ", "),
	)

	if err := smtp.SendMail(c.address, c.auth, from, to, buf.Bytes()); err != nil {
		return errors.ErrInternalServiceOperationFailed.WithError(err, "sending mail failed")
	}

	return nil
}
