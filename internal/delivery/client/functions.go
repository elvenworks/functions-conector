package client

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	apiv2 "cloud.google.com/go/functions/apiv2"
	functionspb "cloud.google.com/go/functions/apiv2/functionspb"
	"cloud.google.com/go/logging"
	"cloud.google.com/go/logging/logadmin"
	"github.com/elvenworks/functions-conector/domain"
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

func (c *Client) GetLastFunctionsRun(config *functions.Config, name, validationString, locations string, seconds time.Duration) (lastRun *domain.FunctionsLastRun, err error) {
	const functionExecutionTook = "Function execution took"
	var entries []*logging.Entry
	var entriesClean []*logging.Entry
	lastSeconds := time.Now().Add(-seconds * time.Second).UTC().Format(time.RFC3339)

	client, err := apiv2.NewFunctionClient(config.Context, config.Option)
	if err != nil {
		return lastRun, fmt.Errorf("failed to create client: %v", err)
	}

	defer client.Close()

	req := &functionspb.GetFunctionRequest{
		Name: fmt.Sprintf(`projects/%s/locations/%s/functions/%s`, config.Credentials.ProjectID, locations, name),
	}

	resp, err := client.GetFunction(config.Context, req)
	if err != nil {
		return lastRun, fmt.Errorf("failed to get function: %v", err)
	}

	if resp.State != functionspb.Function_ACTIVE {
		return lastRun, fmt.Errorf("%v is not activated", name)
	}

	if resp.Environment != functionspb.Environment_GEN_1 {
		return lastRun, fmt.Errorf("%v is %s not %s ", name, resp.Environment.String(), functionspb.Environment_GEN_1.String())
	}

	iter := c.client.Entries(config.Context,
		logadmin.Filter(fmt.Sprintf(`resource.type="cloud_function" AND resource.labels.function_name="%s" AND timestamp >= "%s" AND resource.labels.region="%s"`, name, lastSeconds, locations)),
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

	if validationString != "" {
		match := c.matchReturn(validationString, entriesClean[0].Payload)

		if !match {
			return lastRun, fmt.Errorf("%v - %v", entriesClean[0].Payload, entriesClean[1].Payload)
		}
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

func (c *Client) GetLastFunctionsRunGen2(config *functions.Config, name, locations string, seconds time.Duration) (lastRun *domain.FunctionsLastRun, err error) {
	var entries []*logging.Entry
	var entriesClean []*logging.Entry
	lastSeconds := time.Now().Add(-seconds * time.Second).UTC().Format(time.RFC3339)

	client, err := apiv2.NewFunctionClient(config.Context, config.Option)
	if err != nil {
		return lastRun, fmt.Errorf("failed to create client: %v", err)
	}

	defer client.Close()

	req := &functionspb.GetFunctionRequest{
		Name: fmt.Sprintf(`projects/%s/locations/%s/functions/%s`, config.Credentials.ProjectID, locations, name),
	}

	resp, err := client.GetFunction(config.Context, req)
	if err != nil {
		return lastRun, fmt.Errorf("failed to get function: %v", err)
	}

	if resp.State != functionspb.Function_ACTIVE {
		return lastRun, fmt.Errorf("%v is not activated", name)
	}

	if resp.Environment != functionspb.Environment_GEN_2 {
		return lastRun, fmt.Errorf("%v is %s not %s ", name, resp.Environment.String(), functionspb.Environment_GEN_2.String())
	}

	iter := c.client.Entries(config.Context,
		logadmin.Filter(fmt.Sprintf(`resource.type = "cloud_run_revision" AND resource.labels.service_name = "%s" AND timestamp >= "%s" AND resource.labels.location = "%s"`, name, lastSeconds, locations)),
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

	onError := false
	statusCode := 0

	for _, entry := range entries {
		if entry.HTTPRequest != nil {
			statusCode = entry.HTTPRequest.Status
			if entry.HTTPRequest.Status > 299 {
				entriesClean = append(entriesClean, entry)
				onError = true
				break
			}
			entriesClean = append(entriesClean, entry)
		}
	}

	if len(entriesClean) == 0 {
		return nil, nil
	}

	lastRun = &domain.FunctionsLastRun{
		Timestamp: entriesClean[0].Timestamp,
	}

	if onError {
		return lastRun, fmt.Errorf("Last run returns status code %v", statusCode)
	}

	return lastRun, nil

}
