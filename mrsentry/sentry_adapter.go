package mrsentry

import (
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/mondegor/go-sysmess/mrerr"

	"github.com/mondegor/go-webcore/mrcore"
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

	// Options - опции для создания Adapter.
	Options struct {
		Dsn              string
		Environment      string
		TracesSampleRate float64
		FlushTimeout     time.Duration
		StackTraceBounds []string
		IsDebug          bool
	}
)

// New - создаёт объект Adapter.
func New(opts Options) (*Adapter, error) {
	sentryOpts := sentry.ClientOptions{
		Dsn:              opts.Dsn,
		Environment:      opts.Environment,
		TracesSampleRate: opts.TracesSampleRate,
		Debug:            opts.IsDebug,
	}

	if len(opts.StackTraceBounds) > 0 {
		sentryOpts.BeforeSend = filterStackTrace(opts.StackTraceBounds)
	}

	client, err := sentry.NewClient(sentryOpts)
	if err != nil {
		return nil, mrcore.ErrStorageConnectionFailed.Wrap(err, connectionName)
	}

	if opts.FlushTimeout == 0 {
		opts.FlushTimeout = defaultFlushTimeout
	}

	return &Adapter{
		client:       client,
		flushTimeout: opts.FlushTimeout,
	}, nil
}

// Cli - comment method.
func (a *Adapter) Cli() *sentry.Client {
	return a.client
}

// CaptureAppError - comment method.
func (a *Adapter) CaptureAppError(err *mrerr.AppError) (instanceID string) {
	sentry.CurrentHub().WithScope(
		func(scope *sentry.Scope) {
			// TODO: добавить отправку аргументов и атрибутов ошибки с помощью scope.SetExtras()
			scope.SetTag(errorKindTagName, err.Kind().String())
			scope.AddEventProcessor(
				func(event *sentry.Event, _ *sentry.EventHint) *sentry.Event {
					var unwrappedErr errorInfo = err

					for i := len(event.Exception) - 1; i >= 0; i-- {
						event.Exception[i].Type += " (" + unwrappedErr.Code() + ")"

						unwrappedErr = unwrapErrorInfo(unwrappedErr)

						if unwrappedErr == nil {
							break
						}
					}

					return event
				},
			)

			if eventID := a.client.CaptureException(mrerr.WithoutStackTrace(err), nil, scope); eventID != nil {
				instanceID = string(*eventID)
			}
		},
	)

	return instanceID
}

// Close - comment method.
func (a *Adapter) Close() error {
	if a.client == nil {
		return mrcore.ErrStorageConnectionIsNotOpened.New(connectionName)
	}

	a.client.Flush(a.flushTimeout)
	a.client = nil

	return nil
}
