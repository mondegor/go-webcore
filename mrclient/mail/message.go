package mail

import (
	"encoding/base64"
	"net/mail"
	"net/textproto"
	"strings"

	"github.com/mondegor/go-sysmess/errors"
)

const (
	defaultContentType          = "text/plain"
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

// ErrInternalParsingAddressFailed - parsing address failed.
var ErrInternalParsingAddressFailed = errors.NewInternalProto("parsing address failed")

// NewMessage - создаёт объект Message.
// Где from - электронный адрес отправителя, to - электронный адрес получателя.
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

// Header - метаинформация о сообщении в MIME формате.
func (d *Message) Header() textproto.MIMEHeader {
	return d.header
}

// From - отправитель сообщения.
func (d *Message) From() string {
	return d.from
}

// To - получатели сообщения.
func (d *Message) To() []string {
	return d.to
}

func createMessageHeader(msg *messageOptions, from, to string) textproto.MIMEHeader {
	header := make(textproto.MIMEHeader)

	header.Set("Mime-Version", "1.0")
	header.Set("Subject", encodeValue(msg.subject, "UTF-8"))
	header.Set("Content-Type", msg.contentType+"; charset=\"UTF-8\"")
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

func encodeValue(value, charset string) string {
	return "=?" + charset + "?B?" + base64.StdEncoding.EncodeToString([]byte(value)) + "?="
}
