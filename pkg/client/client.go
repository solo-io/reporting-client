package client

import (
	"context"
	"os"
	"time"

	api "github.com/solo-io/reporting-client/pkg/api/v1"
	"google.golang.org/grpc"
)

const (
	// set this env var to the string "true" to prevent usage from being reported
	DisableUsageVar  = "USAGE_REPORTING_DISABLE"
	errorSendTimeout = time.Second * 10
)

// a type that knows how to load the usage payload you want to report
type UsagePayloadReader interface {
	GetPayload() (map[string]string, error)
}

type ReportingServiceClientBuilder func() (api.ReportingServiceClient, error)

var defaultReportingServiceClientBuilder = func(url string) ReportingServiceClientBuilder {
	return func() (api.ReportingServiceClient, error) {
		clientConn, err := grpc.Dial(url, grpc.WithInsecure())
		if err != nil {
			return nil, err
		}

		client := api.NewReportingServiceClient(clientConn)
		return client, nil
	}
}

//go:generate mockgen -destination mocks/mock_usage_payload_reader.go -package mocks github.com/solo-io/reporting-client/pkg/client UsagePayloadReader
//go:generate mockgen -destination mocks/mock_reporting_service_client.go -package mocks github.com/solo-io/reporting-client/pkg/api/v1 ReportingServiceClient

type Client interface {
	StartReportingUsage(ctx context.Context, interval time.Duration) <-chan error
}

type client struct {
	usagePayloadReader UsagePayloadReader
	usageClientBuilder ReportingServiceClientBuilder
	metadata           *api.InstanceMetadata
}

var _ Client = &client{}

func NewUsageClient(usageServerUrl string, usagePayloadReader UsagePayloadReader, instanceMetadata *api.InstanceMetadata) *client {
	return newUsageClient(
		usagePayloadReader,
		instanceMetadata,
		defaultReportingServiceClientBuilder(usageServerUrl),
	)
}

// visible for testing
func newUsageClient(
	usagePayloadReader UsagePayloadReader,
	instanceMetadata *api.InstanceMetadata,
	reportingServiceClientBuilder ReportingServiceClientBuilder,
) *client {
	return &client{
		usagePayloadReader: usagePayloadReader,
		usageClientBuilder: reportingServiceClientBuilder,
		metadata:           instanceMetadata,
	}
}

func (c *client) StartReportingUsage(ctx context.Context, interval time.Duration) <-chan error {
	errorChan := make(chan error)
	if os.Getenv(DisableUsageVar) == "true" {
		close(errorChan)
		return errorChan
	}

	ticker := time.NewTicker(interval)
	go func() {
		for {
			select {
			case <-ticker.C:
				payload, err := c.usagePayloadReader.GetPayload()
				if err != nil {
					sendWithTimeout(errorChan, ErrorReadingPayload(err))
					continue
				}

				client, err := c.usageClientBuilder()
				if err != nil {
					sendWithTimeout(errorChan, ErrorConnecting(err))
					continue
				}
				_, err = client.ReportUsage(ctx, &api.UsageRequest{
					InstanceMetadata: c.metadata,
					Payload:          payload,
				})
				if err != nil {
					sendWithTimeout(errorChan, ErrorSendingUsage(err))
				}
			case <-ctx.Done():
				close(errorChan)
				return
			}
		}
	}()

	return errorChan
}

// we don't want to block this whole goroutine if no one is listening for errors,
// so if a receiver isn't ready after the timeout, give up and continue
func sendWithTimeout(errorChan chan<- error, err error) {
	select {
	case errorChan <- err:
	case <-time.After(errorSendTimeout):
	}
}
