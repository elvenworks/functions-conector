package functions

import (
	"time"

	"github.com/elvenworks/functions-conector/domain"
	"github.com/elvenworks/functions-conector/internal/delivery/client"
	"github.com/elvenworks/functions-conector/internal/driver/functions"
	"github.com/sirupsen/logrus"
)

type Secret struct {
	JsonCredentials []byte
}

type Functions struct {
	config *functions.Config
}

func InitFunctions(secret Secret) (f *Functions, err error) {
	config, err := functions.NewConfig(secret.JsonCredentials)

	if err != nil {
		logrus.Errorf("Failed loading config, err: %s\n", err)
		return nil, err
	}

	return &Functions{
		config: config,
	}, nil

}

func (f *Functions) GetLastFunctionsRun(name, validationString, locations string, seconds time.Duration) (lastRun *domain.FunctionsLastRun, err error) {

	client, err := client.NewClient(f.config)

	if err != nil {
		return nil, err
	}

	lastRun, err = client.GetLastFunctionsRun(f.config, name, validationString, locations, seconds)

	client.Close()

	return lastRun, err
}
