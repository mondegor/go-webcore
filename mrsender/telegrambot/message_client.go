package telegrambot

import (
	"context"
	"strconv"

	tgclient "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/mondegor/go-sysmess/mrerr/mr"
	"github.com/mondegor/go-sysmess/mrtrace"
)

const (
	telegramBotClientName = "TelegramBotClient"
)

// MessageClient - адаптер для отправки сообщений через telegram.
type (
	MessageClient struct {
		botAPI *tgclient.BotAPI
		tracer mrtrace.Tracer
	}
)

// NewMessageClient - создаёт объект MessageClient.
func NewMessageClient(apiToken string, tracer mrtrace.Tracer) (client *MessageClient, err error) {
	botAPI, err := tgclient.NewBotAPI(apiToken)
	if err != nil {
		return nil, mr.ErrInternal.Wrap(err)
	}

	return &MessageClient{
		botAPI: botAPI,
		tracer: tracer,
	}, nil
}

// SendToChat - отправляет сообщение в указанный чат.
func (c *MessageClient) SendToChat(ctx context.Context, chatKey, message string) error {
	chatID, err := strconv.ParseInt(chatKey, 10, 64)
	if err != nil {
		return mr.ErrUseCaseIncorrectInternalInputData.Wrap(err, "chatKey", chatKey)
	}

	c.tracer.Trace(
		ctx,
		"source", telegramBotClientName,
		"cmd", "SendToChat",
		"ChatKey", chatKey,
		"Message", message,
	)

	msg := tgclient.NewMessage(chatID, message)
	msg.ParseMode = "Markdown"

	if _, err = c.botAPI.Send(msg); err != nil {
		return mr.ErrUseCaseOperationFailed.Wrap(err)
	}

	return nil
}
