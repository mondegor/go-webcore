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

// WithContentType - устанавливает тип содержимого письма (contentType).
// Например: "text/plain" для обычного текста или "text/html" для HTML.
func WithContentType(value string) MessageOption {
	return func(o *messageOptions) {
		o.contentType = value
	}
}

// WithSubject - устанавливает тему письма (subject).
// Если не указана, будет использована тема по умолчанию.
func WithSubject(value string) MessageOption {
	return func(o *messageOptions) {
		o.subject = value
	}
}

// WithUseExtendEmailFormat - разрешает или запрещает использование расширенного формата email.
// В расширенном формате email может содержать имя получателя, например: "Иван Иванов <ivan@example.com>".
func WithUseExtendEmailFormat(value bool) MessageOption {
	return func(o *messageOptions) {
		o.useExtendEmailFormat = value
	}
}

// WithCC - устанавливает список получателей копии письма (carbon copy).
// Параметр value должен содержать email-адреса, разделённые запятыми.
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

// WithReplyTo - устанавливает email для ответов на письмо (Reply-To).
// Если не указан, ответы будут отправляться на адрес отправителя по умолчанию.
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

// WithReturnEmail - устанавливает обратный email для возврата письма (Return-Path).
// Используется почтовыми серверами для обработки недоставленных писем.
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
