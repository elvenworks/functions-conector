package functions

import (
	"context"
	"encoding/json"
	"fmt"

	apiv1 "cloud.google.com/go/functions/apiv1"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"

	"github.com/elvenworks/functions-conector/domain"
)

type Config struct {
	Context     context.Context
	Credentials domain.Credentials
	Option      option.ClientOption
}

func NewConfig(jsonCredentials []byte) (c *Config, err error) {
	var context = context.Background()
	var credentials domain.Credentials

	if err := json.Unmarshal(jsonCredentials, &credentials); err != nil {
		logrus.Errorf("Failed to unmarshal credentials to functions, err: %s\n", err)
		return nil, err
	}

	if credentials.ProjectID == "" {
		return nil, fmt.Errorf("projectID not found to functions")
	}

	creds, err := google.CredentialsFromJSON(context, jsonCredentials, apiv1.DefaultAuthScopes()...)

	if err != nil {
		return nil, err
	}

	config := &Config{Context: context, Credentials: credentials, Option: option.WithCredentials(creds)}

	return config, nil

}
