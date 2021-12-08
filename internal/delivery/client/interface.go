package client

import (
	"time"

	"github.com/elvenworks/functions-conector/internal/domain"
	"github.com/elvenworks/functions-conector/internal/driver/functions"
)

type IFunctionsClient interface {
	GetLastFunctionsRun(config *functions.Config, name, validationString string, seconds time.Duration) (lastRun *domain.FunctionsLastRun, err error)
	matchReturn(expected string, content interface{}) bool
}
