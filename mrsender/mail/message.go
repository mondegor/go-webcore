package mail

import (
	"net/mail"
	"net/textproto"
	"strings"
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

	message struct {
		contentType          string
		subject              string
		useExtendEmailFormat bool
		parser               *mail.AddressParser
		cc                   []*mail.Address
		replyTo              *mail.Address
		returnEmail          string
		err                  error
	}
)

// NewMessage - создаёт объект Message.
func NewMessage(from, to string, opts ...MessageOption) (*Message, error) {
	emailParser := mail.AddressParser{}

	fromEmail, err := emailParser.Parse(from)
	if err != nil {
		return nil, err
	}

	toEmail, err := emailParser.Parse(to)
	if err != nil {
		return nil, err
	}

	wm := message{
		contentType:          defaultContentType,
		subject:              defaultMessageSubject,
		useExtendEmailFormat: defaultUseExtendEmailFormat,
		parser:               &emailParser,
	}

	for _, opt := range opts {
		opt(&wm)

		if wm.err != nil {
			return nil, wm.err
		}
	}

	if !wm.useExtendEmailFormat {
		fromEmail.Name = ""
		toEmail.Name = ""

		for i := range wm.cc {
			wm.cc[i].Name = ""
		}

		if wm.replyTo != nil {
			wm.replyTo.Name = ""
		}
	}

	if wm.returnEmail == "" {
		wm.returnEmail = fromEmail.Address
	}

	header := createMessageHeader(&wm, fromEmail.String(), toEmail.String())

	toList := make([]string, len(wm.cc)+1)
	toList[0] = toEmail.Address

	if len(wm.cc) > 0 {
		for i := range wm.cc {
			toList[i+1] = wm.cc[i].Address
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

func createMessageHeader(m *message, from, to string) textproto.MIMEHeader {
	header := make(textproto.MIMEHeader)

	header.Set("Mime-Version", "1.0")
	header.Set("Subject", m.subject)
	header.Set("Content-Type", m.contentType+"; charset=\""+"UTF-8"+"\"")
	header.Set("From", from)
	header.Set("To", to)

	if len(m.cc) > 0 {
		var buf strings.Builder

		buf.WriteString(m.cc[0].String())

		for i := 1; i < len(m.cc); i++ {
			buf.WriteString(", ")
			buf.WriteString(m.cc[i].String())
		}

		header.Set("cc", buf.String())
	}

	if m.replyTo != nil {
		header.Set("Reply-To", m.replyTo.String())
	}

	header.Set("Return-Path", m.returnEmail)

	return header
}
