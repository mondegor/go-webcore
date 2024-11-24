package main

import (
	"context"

	"github.com/mondegor/go-webcore/mrlog"
	"github.com/mondegor/go-webcore/mrsender/mail"
	"github.com/mondegor/go-webcore/mrsender/mail/smtp"
)

func main() {
	logger := mrlog.New(mrlog.TraceLevel).With().Str("service", "smtpmail").Logger()
	ctx := mrlog.WithContext(context.Background(), logger)

	smtpHost := "{host}" // smtp.gmail.com
	smtpPort := "{port}" // 587 с поддержкой STARTTLS
	smtpUsername := "{user_login}"
	smtpPassword := "{user_password}"

	from := "Test Sender <{email_from}>" // from_test@gmail.com
	to := "Test Recipient <{email_to}>"  // to_test@gmail.com
	body := "The Test Content"

	msg, err := mail.NewMessage(
		from,
		to,
		mail.WithSubject("The Test Subject"),
		// mail.WithCC("Test Recipient2 <{email_to2}>, Test Recipient3 <{email_to3}>"),
		// mail.WithUseExtendEmailFormat(false),
	)
	if err != nil {
		logger.Fatal().Err(err).Send()
	}

	mailClient := smtp.NewMailClient(smtpHost, smtpPort, smtpUsername, smtpPassword)

	if err = mailClient.SendMail(ctx, msg.From(), msg.To(), msg.Header(), body); err != nil {
		logger.Fatal().Err(err).Send()
	}
}
