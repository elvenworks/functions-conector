package client

import (
	"context"
	"time"

	"cloud.google.com/go/logging/logadmin"
	"github.com/elvenworks/functions-conector/internal/domain"
	"github.com/elvenworks/functions-conector/internal/driver/functions"
)

type ILoggingClient interface {
	Close() error
	DeleteLog(ctx context.Context, logID string) error
	Entries(ctx context.Context, opts ...logadmin.EntriesOption) *logadmin.EntryIterator
	Logs(ctx context.Context) *logadmin.LogIterator
}

type IFunctionsClient interface {
	GetLastFunctionsRun(config *functions.Config, name, validationString string, seconds time.Duration) (lastRun *domain.FunctionsLastRun, err error)
	Close()
}
