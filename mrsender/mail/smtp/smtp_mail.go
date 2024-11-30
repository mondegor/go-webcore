package smtp

import (
	"bytes"
	"context"
	"net/smtp"
	"net/textproto"
	"strings"

	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrlog"
)

const (
	smtpMailClientName = "SmtpMailClient"
	defaultMailSubject = "The mail without a subject"
)

type (
	// MailClient - адаптер для отправки электронных писем через SMTP.
	MailClient struct {
		address string
		auth    smtp.Auth
	}
)

// NewMailClient - создаёт объект MailClient.
func NewMailClient(host, port, username, password string) *MailClient {
	return &MailClient{
		address: host + ":" + port,
		auth:    smtp.PlainAuth("", username, password, host),
	}
}

// SendMail - отправляет электронное письмо указанному адресату.
// Где from - электронный адрес отправителя, to - электронные адреса получателей.
func (c *MailClient) SendMail(ctx context.Context, from string, to []string, header textproto.MIMEHeader, body string) error {
	if from == "" {
		return mrcore.ErrUseCaseRequiredDataIsEmpty.New("From address")
	}

	if len(to) == 0 || to[0] == "" {
		return mrcore.ErrUseCaseRequiredDataIsEmpty.New("To address")
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
		buf.WriteString(defaultMailSubject)
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

	mrlog.Ctx(ctx).
		Trace().
		Str("source", smtpMailClientName).
		Str("cmd", "SendMail").
		MsgFunc(
			func() string {
				return "SMTP-Header: \n" + buf.String() + "\n" +
					"SMTP-From: " + from + "\n" +
					"SMTP-To: " + strings.Join(to, ", ")
			},
		)

	if err := smtp.SendMail(c.address, c.auth, from, to, buf.Bytes()); err != nil {
		return mrcore.ErrUseCaseOperationFailed.Wrap(err)
	}

	return nil
}
