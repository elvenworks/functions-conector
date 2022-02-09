package functions

import (
	"time"

	"github.com/elvenworks/functions-conector/domain"
	"github.com/stretchr/testify/mock"
)

type MockFunctions struct {
	mock.Mock
}

func (mock MockFunctions) GetLastFunctionsRun(name, validationString, locations string, seconds time.Duration) (lastRun *domain.FunctionsLastRun, err error) {
	args := mock.Called(name, validationString, locations, seconds)
	return args.Get(0).(*domain.FunctionsLastRun), args.Error(1)
}
