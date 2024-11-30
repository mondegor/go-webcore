package mail

type (
	// MessageOption - настройка объекта Message.
	MessageOption func(m *message)
)

// WithContentType - устанавливает опцию contentType для Message.
func WithContentType(value string) MessageOption {
	return func(m *message) {
		if m.contentType != "" {
			m.contentType = value
		}
	}
}

// WithSubject - устанавливает тему сообщения.
func WithSubject(value string) MessageOption {
	return func(m *message) {
		m.subject = value
	}
}

// WithUseExtendEmailFormat - устанавливает опцию, позволяющую использовать расширенный формат электронного адреса.
func WithUseExtendEmailFormat(value bool) MessageOption {
	return func(m *message) {
		m.useExtendEmailFormat = value
	}
}

// WithCC - устанавливает список получателей копии письма разделённых через ",".
func WithCC(value string) MessageOption {
	return func(m *message) {
		list, err := m.parser.ParseList(value)
		if err != nil {
			m.err = err

			return
		}

		m.cc = list
	}
}

// WithReplyTo - устанавливает электронный адрес по умолчанию при ответе на письмо.
func WithReplyTo(value string) MessageOption {
	return func(m *message) {
		email, err := m.parser.Parse(value)
		if err != nil {
			m.err = err

			return
		}

		m.replyTo = email
	}
}

// WithReturnEmail - устанавливает опцию обратного электронного адреса (служебный).
func WithReturnEmail(value string) MessageOption {
	return func(m *message) {
		email, err := m.parser.Parse(value)
		if err != nil {
			m.err = err

			return
		}

		m.returnEmail = email.Address
	}
}
