package client

import (
	"context"
	"time"

	"cloud.google.com/go/logging/logadmin"
	"github.com/elvenworks/functions-conector/internal/driver/functions"
	"github.com/stretchr/testify/mock"
)

type MockClient struct {
	mock.Mock
}

func (mock MockClient) GetLastFunctionsRun(config *functions.Config, name, validationString string, seconds time.Duration) (err error) {
	args := mock.Called(config, name, validationString, seconds)
	return args.Error(0)
}

func (mock MockClient) Close() {
}

func (mock MockClient) matchReturn(expected string, content interface{}) bool {
	args := mock.Called()
	return args.Bool(0)
}

type MockLogging struct {
	mock.Mock
}

func (mock MockLogging) Close() error {
	args := mock.Called()
	return args.Error(0)
}

func (mock MockLogging) DeleteLog(ctx context.Context, logID string) error {
	args := mock.Called(ctx, logID)
	return args.Error(0)
}

func (mock MockLogging) Entries(ctx context.Context, opts ...logadmin.EntriesOption) *logadmin.EntryIterator {
	args := mock.Called(ctx, opts)
	return args.Get(0).(*logadmin.EntryIterator)
}

func (mock MockLogging) Logs(ctx context.Context) *logadmin.LogIterator {
	args := mock.Called(ctx)
	return args.Get(0).(*logadmin.LogIterator)
}
