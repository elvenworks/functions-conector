package client

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"cloud.google.com/go/logging"
	"cloud.google.com/go/logging/logadmin"
	"github.com/elvenworks/functions-conector/internal/domain"
	"github.com/elvenworks/functions-conector/internal/driver/functions"
	"github.com/sirupsen/logrus"
	"google.golang.org/api/iterator"
)

type Client struct {
	client ILoggingClient
}

func NewClient(config *functions.Config) (f *Client, err error) {
	adminClient, err := logadmin.NewClient(config.Context, config.Credentials.ProjectID, config.Option)

	if err != nil {
		logrus.Errorf("Failed to create logadmin client: %s\n", err)
		return nil, err
	}

	return &Client{
		client: adminClient,
	}, nil
}

func (c *Client) GetLastFunctionsRun(config *functions.Config, name, validationString string, seconds time.Duration) (lastRun *domain.FunctionsLastRun, err error) {
	const functionExecutionTook = "Function execution took"
	var entries []*logging.Entry
	var entriesClean []*logging.Entry
	lastSeconds := time.Now().Add(-seconds * time.Second).Format(time.RFC3339)

	iter := c.client.Entries(config.Context,
		logadmin.Filter(fmt.Sprintf(`resource.type="cloud_function" AND resource.labels.function_name="%s" AND timestamp >= "%s"`, name, lastSeconds)),
		logadmin.NewestFirst(),
	)

	for {
		entry, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			logrus.Errorf("could not read time series value, err: %s\n", err)
			break
		}

		entries = append(entries, entry)
	}

	getLast := false

	for _, entry := range entries {
		if c.matchReturn(functionExecutionTook, entry.Payload) && !getLast {
			entriesClean = append(entriesClean, entry)
			getLast = true
		}

		if !c.matchReturn(functionExecutionTook, entry.Payload) && getLast {
			entriesClean = append(entriesClean, entry)
			break
		}
	}

	if len(entriesClean) == 0 {
		return nil, nil
	}

	lastRun = &domain.FunctionsLastRun{
		Timestamp: entriesClean[0].Timestamp,
	}

	match := c.matchReturn(validationString, entriesClean[0].Payload)

	if !match {
		return lastRun, fmt.Errorf("-> %v - %v", entriesClean[0].Payload, entriesClean[1].Payload)
	}

	return lastRun, nil

}

func (c *Client) Close() {
	defer c.client.Close()
}

func (c *Client) matchReturn(expected string, content interface{}) bool {
	_, err := regexp.Compile(expected)
	if err != nil {
		return strings.Contains(expected, fmt.Sprintf("%v", content))
	} else {
		re := regexp.MustCompile(expected)
		return re.Match([]byte(fmt.Sprintf("%v", content))) || strings.Contains(expected, fmt.Sprintf("%v", content))
	}
}
