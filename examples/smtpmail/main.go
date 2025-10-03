package main

import (
	"context"
	"os"

	"github.com/mondegor/go-sysmess/mrlog"
	"github.com/mondegor/go-sysmess/mrlog/slog"
	"github.com/mondegor/go-sysmess/mrtrace/logtracer"

	"github.com/mondegor/go-webcore/mrsender/mail"
	"github.com/mondegor/go-webcore/mrsender/mail/smtp"
)

func main() {
	logger, _ := slog.NewLoggerAdapter(
		slog.WithWriter(os.Stdout),
	)

	tracer := logtracer.New(logger)

	smtpHost := "{host}" // smtp.gmail.com
	smtpPort := "{port}" // 587 с поддержкой STARTTLS
	smtpUsername := "{user_login}"
	smtpPassword := "{user_password}"

	from := "Test Sender <from_test@gmail.com>"
	to := "Test Recipient <to_test@gmail.com>"
	body := "The Test Content"

	msg, err := mail.NewMessage(
		from,
		to,
		mail.WithSubject("The Test Subject"),
		// mail.WithCC("Test Recipient2 <{email_to2}>, Test Recipient3 <{email_to3}>"),
		// mail.WithUseExtendEmailFormat(false),
	)
	if err != nil {
		mrlog.Fatal(logger, "this is error", "error", err)
	}

	mailClient := smtp.NewMailClient(smtpHost, smtpPort, smtpUsername, smtpPassword, tracer)

	if err = mailClient.SendMail(context.Background(), msg.From(), msg.To(), msg.Header(), body); err != nil {
		os.Exit(1)
	}
}
