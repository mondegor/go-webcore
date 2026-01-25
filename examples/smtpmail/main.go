package main

import (
	"context"
	"os"

	"github.com/mondegor/go-sysmess/mrlog"
	"github.com/mondegor/go-sysmess/mrlog/slog"
	wireslog "github.com/mondegor/go-sysmess/wire/slog"

	"github.com/mondegor/go-webcore/mrclient/mail"
)

func main() {
	logger, _ := slog.NewLoggerAdapter(
		slog.WithWriter(os.Stdout),
	)

	tracer := wireslog.InitTracer(logger)

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

	smtpClient := mail.NewSMTPClient(smtpHost, smtpPort, smtpUsername, smtpPassword, tracer)

	if err = smtpClient.SendMail(context.Background(), msg.From(), msg.To(), msg.Header(), body); err != nil {
		mrlog.Fatal(logger, "this is error", "error", err)
	}
}
