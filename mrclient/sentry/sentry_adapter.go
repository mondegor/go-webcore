package sentry

import (
	"context"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/mondegor/go-sysmess/errors"
	"github.com/mondegor/go-sysmess/errors/kind"
)

// go get -u github.com/getsentry/sentry-go

const (
	connectionName      = "Sentry"
	errorKindTagName    = "error_kind"
	defaultFlushTimeout = 2 * time.Second
)

type (
	// Adapter - comment struct.
	Adapter struct {
		client       *sentry.Client
		flushTimeout time.Duration
	}

	runtimeError interface {
		error

		Kind() kind.Enum
		// Attrs() []any
		Hint() any
	}

	errorHint interface {
		ErrorID() string
		StackTraceIterator() func() (index int, name, file string, line int)
	}
)

// New - создаёт объект Adapter.
func New(dsn string, opts ...Option) (*Adapter, error) {
	o := options{
		flushTimeout: defaultFlushTimeout,
	}

	for _, opt := range opts {
		opt(&o)
	}

	o.sentryOpts.Dsn = dsn

	client, err := sentry.NewClient(o.sentryOpts)
	if err != nil {
		return nil, errors.ErrSystemStorageConnectionFailed.Wrap(err, "source", connectionName)
	}

	return &Adapter{
		client:       client,
		flushTimeout: o.flushTimeout,
	}, nil
}

// Cli - comment method.
func (a *Adapter) Cli() *sentry.Client {
	return a.client
}

// CaptureError - comment method.
func (a *Adapter) CaptureError(_ context.Context, err error) (eventID string) {
	sentry.CurrentHub().WithScope(
		func(scope *sentry.Scope) {
			// TODO: добавить отправку атрибутов из ctx с помощью scope.SetExtras()
			e, ok := err.(runtimeError) //nolint:errorlint
			if !ok {                    // если это не runtime ошибка
				if id := a.client.CaptureException(err, nil, scope); id != nil {
					eventID = string(*id)
				}

				return
			}

			// TODO: добавить отправку аргументов и атрибутов ошибки из ctx с помощью scope.SetExtras()
			scope.SetTag(errorKindTagName, e.Kind().String())

			event := sentry.NewEvent()
			event.Level = sentry.LevelError

			if bag, ok := e.Hint().(errorHint); ok {
				// TODO: stack = strings.Join(stacktrace.ToStrings(bag.StackTraceIterator()), " | ") // TODO: disable function name of stack on prod
				event.EventID = sentry.EventID(bag.ErrorID())
			}

			if id := a.client.CaptureEvent(event, nil, scope); id != nil {
				eventID = string(*id)
			}
		},
	)

	return eventID
}

// Close - comment method.
func (a *Adapter) Close() error {
	if a.client == nil {
		return errors.ErrInternalStorageConnectionIsNotOpened.New("source", connectionName)
	}

	a.client.Flush(a.flushTimeout)
	a.client = nil

	return nil
}
