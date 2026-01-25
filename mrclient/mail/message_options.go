package mail

import "net/mail"

type (
	// MessageOption - настройка объекта Message.
	MessageOption func(o *messageOptions)

	messageOptions struct {
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

// WithContentType - устанавливает опцию contentType для Message.
func WithContentType(value string) MessageOption {
	return func(o *messageOptions) {
		o.contentType = value
	}
}

// WithSubject - устанавливает тему сообщения.
func WithSubject(value string) MessageOption {
	return func(o *messageOptions) {
		o.subject = value
	}
}

// WithUseExtendEmailFormat - устанавливает опцию, позволяющую использовать расширенный формат электронного адреса.
func WithUseExtendEmailFormat(value bool) MessageOption {
	return func(o *messageOptions) {
		o.useExtendEmailFormat = value
	}
}

// WithCC - устанавливает список получателей копии письма разделённых через ",".
func WithCC(value string) MessageOption {
	return func(o *messageOptions) {
		list, err := o.parser.ParseList(value)
		if err != nil {
			o.err = err

			return
		}

		o.cc = list
	}
}

// WithReplyTo - устанавливает электронный адрес по умолчанию при ответе на письмо.
func WithReplyTo(value string) MessageOption {
	return func(o *messageOptions) {
		email, err := o.parser.Parse(value)
		if err != nil {
			o.err = err

			return
		}

		o.replyTo = email
	}
}

// WithReturnEmail - устанавливает опцию обратного электронного адреса (служебный).
func WithReturnEmail(value string) MessageOption {
	return func(o *messageOptions) {
		email, err := o.parser.Parse(value)
		if err != nil {
			o.err = err

			return
		}

		o.returnEmail = email.Address
	}
}
