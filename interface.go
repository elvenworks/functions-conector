package functions

import (
	"time"

	"github.com/elvenworks/functions-conector/domain"
)

type IFunctions interface {
	GetLastFunctionsRun(name, validationString, locations string, seconds time.Duration) (lastRun *domain.FunctionsLastRun, err error)
}
