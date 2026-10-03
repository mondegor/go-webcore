package mail

import (
	"encoding/base64"
	"mime"
	"net/mail"
	"net/textproto"
	"strings"

	"github.com/mondegor/go-core/errors"
)

const (
	defaultContentType          = "text/plain"
	messageCharset              = "UTF-8"
	defaultMessageSubject       = "The mail without a subject"
	defaultUseExtendEmailFormat = true
)

type (
	// Message - подготовленное сообщения для
	// его отправки в виде электронного письма.
	Message struct {
		header textproto.MIMEHeader
		from   string
		to     []string
	}
)

var (
	// ErrInternalParsingAddressFailed - ошибка при неудачном разборе email-адреса.
	ErrInternalParsingAddressFailed = errors.NewInternalProto("parsing address failed")

	// ErrInternalParsingContentTypeFailed - ошибка при неудачном разборе типа содержимого письма.
	ErrInternalParsingContentTypeFailed = errors.NewInternalProto("parsing content type failed")
)

// NewMessage - создаёт и настраивает объект Message для отправки электронного письма.
// Параметры:
//   - from - email отправителя;
//   - to - email основного получателя;
//
// Опциональные параметры opts позволяют настроить тему, тип контента, копии и другие заголовки.
// Возвращает ошибку ErrInternalParsingAddressFailed при некорректном формате email
// и ErrInternalParsingContentTypeFailed при некорректном типе содержимого.
func NewMessage(from, to string, opts ...MessageOption) (*Message, error) {
	emailParser := mail.AddressParser{}

	fromEmail, err := emailParser.Parse(from)
	if err != nil {
		return nil, ErrInternalParsingAddressFailed.Wrap(err)
	}

	toEmail, err := emailParser.Parse(to)
	if err != nil {
		return nil, ErrInternalParsingAddressFailed.Wrap(err)
	}

	o := messageOptions{
		subject:              defaultMessageSubject,
		useExtendEmailFormat: defaultUseExtendEmailFormat,
		parser:               &emailParser,
	}

	for _, opt := range opts {
		opt(&o)

		if o.err != nil {
			return nil, o.err
		}
	}

	if o.contentType == "" {
		o.contentType = defaultContentType
	}

	contentType, err := normalizeContentType(o.contentType)
	if err != nil {
		return nil, err
	}

	o.contentType = contentType

	if !o.useExtendEmailFormat {
		fromEmail.Name = ""
		toEmail.Name = ""

		for i := range o.cc {
			o.cc[i].Name = ""
		}

		if o.replyTo != nil {
			o.replyTo.Name = ""
		}
	}

	if o.returnEmail == "" {
		o.returnEmail = fromEmail.Address
	}

	header := createMessageHeader(&o, fromEmail.String(), toEmail.String())

	toList := make([]string, len(o.cc)+1)
	toList[0] = toEmail.Address

	if len(o.cc) > 0 {
		for i := range o.cc {
			toList[i+1] = o.cc[i].Address
		}
	}

	return &Message{
		header: header,
		from:   fromEmail.Address,
		to:     toList,
	}, nil
}

// Header - возвращает MIME-заголовки сообщения.
// Содержит метаданные письма: тему, отправителя, получателя, тип контента и т.д.
func (d *Message) Header() textproto.MIMEHeader {
	return d.header
}

// From - возвращает email-адрес отправителя письма.
func (d *Message) From() string {
	return d.from
}

// To - возвращает список email-адресов получателей письма.
// Включает основного получателя и всех получателей копии (CC).
func (d *Message) To() []string {
	return d.to
}

func createMessageHeader(msg *messageOptions, from, to string) textproto.MIMEHeader {
	header := make(textproto.MIMEHeader)

	header.Set("Mime-Version", "1.0")
	header.Set("Subject", encodeValue(msg.subject, messageCharset))
	header.Set("Content-Type", msg.contentType)
	header.Set("From", from)
	header.Set("To", to)

	if len(msg.cc) > 0 {
		var buf strings.Builder

		buf.WriteString(msg.cc[0].String())

		for i := 1; i < len(msg.cc); i++ {
			buf.WriteString(", ")
			buf.WriteString(msg.cc[i].String())
		}

		header.Set("cc", buf.String())
	}

	if msg.replyTo != nil {
		header.Set("Reply-To", msg.replyTo.String())
	}

	header.Set("Return-Path", msg.returnEmail)

	return header
}

// normalizeContentType - разбирает тип содержимого и собирает его заново:
// медиа-тип приводится к нижнему регистру, для text/* charset допускается только UTF-8
// (тело и тема письма всегда в UTF-8): отсутствующий выставляется, другой приводит к ошибке;
// прочие параметры сохраняются.
func normalizeContentType(value string) (string, error) {
	mediaType, params, err := mime.ParseMediaType(value)
	if err != nil {
		return "", ErrInternalParsingContentTypeFailed.Wrap(err, "contentType", value)
	}

	// ParseMediaType допускает тип без подтипа (например, "text") и wildcard-типы
	// (например, "text/*", "*/*"), для письма они некорректны.
	if !strings.Contains(mediaType, "/") || strings.Contains(mediaType, "*") {
		return "", ErrInternalParsingContentTypeFailed.New("contentType", value)
	}

	if strings.HasPrefix(mediaType, "text/") {
		if charset, ok := params["charset"]; ok && !strings.EqualFold(charset, messageCharset) {
			return "", ErrInternalParsingContentTypeFailed.New("contentType", value)
		}

		params["charset"] = messageCharset
	}

	return mime.FormatMediaType(mediaType, params), nil
}

func encodeValue(value, charset string) string {
	return "=?" + charset + "?B?" + base64.StdEncoding.EncodeToString([]byte(value)) + "?="
}
