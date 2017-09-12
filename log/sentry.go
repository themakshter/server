package log

import (
	raven "github.com/getsentry/raven-go"
)

type sentry struct {
	client *raven.Client
}

// NewSentryErrorTracker returns an ErrorTracker backed by Sentry
func NewSentryErrorTracker(dsn string) (ErrorTracker, error) {
	client, err := raven.New(dsn)
	if err != nil {
		return nil, err
	}
	Info("Sentry error logging configured", nil)
	return &sentry{
		client: client,
	}, nil
}

func (s *sentry) Name() string {
	return "sentry"
}

func (s *sentry) TrackError(err error, tags map[string]string) {
	s.client.CaptureError(err, tags)
}

func (s *sentry) TrackFatal(err error, tags map[string]string) {
	s.client.CaptureErrorAndWait(err, tags)
}
