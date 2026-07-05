package telegram

import (
	"context"
	"strconv"

	tgclient "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/mondegor/go-core/errors"
	"github.com/mondegor/go-core/mrtrace"
)

const (
	botClientName = "TelegramBotClient"
)

type (
	// BotClient - адаптер для отправки сообщений через telegram.
	BotClient struct {
		botAPI *tgclient.BotAPI
		tracer mrtrace.Tracer
	}
)

// NewBotClient - создаёт и настраивает BotClient для работы с Telegram Bot API.
// Параметры:
//   - apiToken - токен авторизации Telegram бота (получается от @BotFather);
//   - tracer - используется для трассировки операций отправки сообщений;
func NewBotClient(
	apiToken string,
	tracer mrtrace.Tracer,
) (*BotClient, error) {
	botAPI, err := tgclient.NewBotAPI(apiToken)
	if err != nil {
		return nil, errors.WrapInternalError(err, "creating telegram botAPI failed")
	}

	return &BotClient{
		botAPI: botAPI,
		tracer: tracer,
	}, nil
}

// SendToChat - отправляет текстовое сообщение в указанный Telegram-чат.
// Параметры:
//   - chatKey - числовой идентификатор чата (chat_id) в виде строки;
//   - message - текст сообщения с поддержкой Markdown-разметки;
func (c *BotClient) SendToChat(ctx context.Context, chatKey, message string) error {
	chatID, err := strconv.ParseInt(chatKey, 10, 64)
	if err != nil {
		return errors.ErrInternalIncorrectInputData.Wrap(err, "chatKey", chatKey)
	}

	c.tracer.Trace(
		ctx,
		"source", botClientName,
		"cmd", "SendToChat",
		"ChatKey", chatKey,
		"Message", message,
	)

	msg := tgclient.NewMessage(chatID, message)
	msg.ParseMode = "Markdown"

	if _, err = c.botAPI.Send(msg); err != nil {
		return errors.ErrInternalServiceOperationFailed.WithError(err, "sending telegram botAPI failed")
	}

	return nil
}
