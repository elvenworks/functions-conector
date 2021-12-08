package functions

import (
	"time"

	"github.com/stretchr/testify/mock"
)

type MockFunctions struct {
	mock.Mock
}

func (mock MockFunctions) GetLastFunctionsRun(name, validationString string, seconds time.Duration) (err error) {
	args := mock.Called(name, validationString, seconds)
	return args.Error(0)
}
