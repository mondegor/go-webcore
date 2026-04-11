package sentry

import (
	"context"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/mondegor/go-sysmess/errors"
	"github.com/mondegor/go-sysmess/errors/runtime/hint"
)

// go get -u github.com/getsentry/sentry-go

const (
	connectionName      = "Sentry"
	errorKindTagName    = "error_kind"
	defaultFlushTimeout = 2 * time.Second
)

type (
	// Adapter - отправляет ошибки и события в систему мониторинга Sentry.
	Adapter struct {
		client       *sentry.Client
		flushTimeout time.Duration
	}
)

// New - создаёт и настраивает объект Adapter для работы с Sentry.
// Возвращает ошибку типа ErrSystemStorageConnectionFailed при неудачном подключении.
func New(dsn string, opts ...Option) (*Adapter, error) {
	o := options{
		adapter: &Adapter{
			flushTimeout: defaultFlushTimeout,
		},
	}

	for _, opt := range opts {
		opt(&o)
	}

	o.sentryOpts.Dsn = dsn

	client, err := sentry.NewClient(o.sentryOpts)
	if err != nil {
		return nil, errors.ErrSystemStorageConnectionFailed.Wrap(err, "source", connectionName)
	}

	o.adapter.client = client

	return o.adapter, nil
}

// Cli - возвращает внутренний клиент Sentry для прямого доступа к расширенным функциям.
func (a *Adapter) Cli() *sentry.Client {
	return a.client
}

// CaptureError - отправляет ошибку в Sentry для мониторинга и логирования.
// Извлекает метаданные ошибки (ErrorID, ErrorKind, StackTrace) с помощью hint.Extract.
// Если ошибка незнакомая (ErrorID пустой), отправляет как исключение.
// Если ошибка знакомая (ErrorID установлен), отправляет как событие с тегами.
// Возвращает идентификатор события (eventID) в Sentry или пустую строку при ошибке.
func (a *Adapter) CaptureError(_ context.Context, err error) (eventID string) {
	sentry.CurrentHub().WithScope(
		func(scope *sentry.Scope) {
			errorHint := hint.Extract(err)

			// TODO: добавить отправку атрибутов из ctx с помощью scope.SetExtras()
			if errorHint.ErrorID() == "" { // если это незнакомая ошибка
				if id := a.client.CaptureException(err, nil, scope); id != nil {
					eventID = string(*id)
				}

				return
			}

			// TODO: добавить отправку аргументов и атрибутов ошибки из ctx с помощью scope.SetExtras()
			scope.SetTag(errorKindTagName, errorHint.ErrorKind().String())

			event := sentry.NewEvent()
			event.Level = sentry.LevelError

			// TODO: stack := strings.Join(stacktrace.ToStrings(errorHint.StackTraceIterator()), " | ") // TODO: disable function name of stack on prod
			event.EventID = sentry.EventID(errorHint.ErrorID())

			if id := a.client.CaptureEvent(event, nil, scope); id != nil {
				eventID = string(*id)
			}
		},
	)

	return eventID
}

// Close - закрывает соединение с Sentry, предварительно отправив все ожидающие события.
// После закрытия объект Adapter становится непригодным для дальнейшего использования.
// Возвращает ошибку типа ErrInternalStorageConnectionIsNotOpened, если клиент не был инициализирован.
func (a *Adapter) Close() error {
	if a.client == nil {
		return errors.ErrInternalStorageConnectionIsNotOpened.New("source", connectionName)
	}

	a.client.Flush(a.flushTimeout)
	a.client = nil

	return nil
}
