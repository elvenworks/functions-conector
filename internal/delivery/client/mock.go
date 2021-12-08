package client

import (
	"time"

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

func (mock MockClient) matchReturn(expected string, content interface{}) bool {
	args := mock.Called()
	return args.Bool(0)
}
