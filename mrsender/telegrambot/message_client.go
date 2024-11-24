package telegrambot

import (
	"context"
	"strconv"

	tgclient "github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrlog"
)

const (
	telegramBotClientName = "TelegramBotClient"
)

// MessageClient - адаптер для отправки сообщений через telegram.
type MessageClient struct {
	botAPI *tgclient.BotAPI
}

// NewMessageClient - создаёт объект MessageClient.
func NewMessageClient(apiToken string) (client *MessageClient, err error) {
	botAPI, err := tgclient.NewBotAPI(apiToken)
	if err != nil {
		return nil, mrcore.ErrInternal.Wrap(err)
	}

	return &MessageClient{
		botAPI: botAPI,
	}, nil
}

// SendToChat - отправляет сообщение в указанный чат.
func (c *MessageClient) SendToChat(ctx context.Context, chatKey, message string) error {
	chatID, err := strconv.ParseInt(chatKey, 10, 64)
	if err != nil {
		return mrcore.ErrUseCaseIncorrectInputData.Wrap(err, "chatID", chatID)
	}

	mrlog.Ctx(ctx).
		Trace().
		Str("source", telegramBotClientName).
		Str("cmd", "SendToChat").
		MsgFunc(
			func() string {
				return "ChatKey: " + chatKey + "\n" +
					"Message: " + message
			},
		)

	msg := tgclient.NewMessage(chatID, message)
	msg.ParseMode = "Markdown"

	if _, err = c.botAPI.Send(msg); err != nil {
		return mrcore.ErrUseCaseOperationFailed.Wrap(err)
	}

	return nil
}
