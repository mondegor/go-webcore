package mrsentry

import (
	"context"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/mondegor/go-sysmess/mrerr/mr"
	"github.com/mondegor/go-sysmess/mrerrors"
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
		DSN              string
		Environment      string
		AppVersion       string
		TracesSampleRate float64
		FlushTimeout     time.Duration
		StackTraceBounds []string
		IsDebug          bool
	}
)

// New - создаёт объект Adapter.
func New(opts Options) (*Adapter, error) {
	sentryOpts := sentry.ClientOptions{
		Dsn:              opts.DSN,
		Environment:      opts.Environment,
		Release:          opts.AppVersion,
		TracesSampleRate: opts.TracesSampleRate,
		Debug:            opts.IsDebug,
	}

	if len(opts.StackTraceBounds) > 0 {
		sentryOpts.BeforeSend = filterStackTrace(opts.StackTraceBounds)
	}

	client, err := sentry.NewClient(sentryOpts)
	if err != nil {
		return nil, mr.ErrStorageConnectionFailed.Wrap(err, connectionName)
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
func (a *Adapter) CaptureAppError(_ context.Context, err *mrerrors.InstantError) (instanceID string) {
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

			// ??????????????????????????????????????????? CastLessVerboseError
			if eventID := a.client.CaptureException(mrerrors.CastLessVerboseError(err), nil, scope); eventID != nil {
				instanceID = string(*eventID)
			}
		},
	)

	return instanceID
}

// Close - comment method.
func (a *Adapter) Close() error {
	if a.client == nil {
		return mr.ErrStorageConnectionIsNotOpened.New(connectionName)
	}

	a.client.Flush(a.flushTimeout)
	a.client = nil

	return nil
}
