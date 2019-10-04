package client

import (
	"context"
	"os"
	"time"

	"github.com/solo-io/go-utils/contextutils"
	api "github.com/solo-io/reporting-client/pkg/api/v1"
	"google.golang.org/grpc"
)

const (
	// set this env var to the string "true" to prevent usage from being reported
	DisableUsageVar = "USAGE_REPORTING_DISABLE"
)

// a type that knows how to load the usage payload you want to report
type UsagePayloadReader interface {
	GetPayload() (map[string]string, error)
}

type ReportingServiceClientBuilder interface {
	BuildClient() (api.ReportingServiceClient, error)
}

type defaultReportingServiceClientBuilder struct {
	url string
}

var _ ReportingServiceClientBuilder = &defaultReportingServiceClientBuilder{}

func (d *defaultReportingServiceClientBuilder) BuildClient() (api.ReportingServiceClient, error) {
	clientConn, err := grpc.Dial(d.url, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	client := api.NewReportingServiceClient(clientConn)
	return client, nil
}

//go:generate mockgen -destination mocks/mock_usage_payload_reader.go -package mocks github.com/solo-io/reporting-client/pkg/client UsagePayloadReader
//go:generate mockgen -destination mocks/mock_reporting_service_client.go -package mocks github.com/solo-io/reporting-client/pkg/api/v1 ReportingServiceClient

type Client interface {
	StartReportingUsage(ctx context.Context, interval time.Duration)
}

type client struct {
	usagePayloadReader UsagePayloadReader
	usageClientBuilder ReportingServiceClientBuilder
	metadata           *api.InstanceMetadata
}

var _ Client = &client{}

// initializes a connection to the grpc server
// returns an error if it is unable to dial the server
func NewUsageClient(usageServerUrl string, usagePayloadReader UsagePayloadReader, instanceMetadata *api.InstanceMetadata) (*client, error) {
	return newUsageClient(
		usagePayloadReader,
		instanceMetadata,
		&defaultReportingServiceClientBuilder{url: usageServerUrl},
	)
}

// visible for testing
func newUsageClient(
	usagePayloadReader UsagePayloadReader,
	instanceMetadata *api.InstanceMetadata,
	reportingServiceClientBuilder ReportingServiceClientBuilder,
) (*client, error) {
	return &client{
		usagePayloadReader: usagePayloadReader,
		usageClientBuilder: reportingServiceClientBuilder,
		metadata:           instanceMetadata,
	}, nil
}

func (c *client) StartReportingUsage(ctx context.Context, interval time.Duration) {
	if os.Getenv(DisableUsageVar) == "true" {
		return
	}

	ticker := time.NewTicker(interval)
	go func() {
		for {
			select {
			case <-ticker.C:
				payload, err := c.usagePayloadReader.GetPayload()
				if err != nil {
					contextutils.LoggerFrom(ctx).Errorf("Encountered error while reading payload: %s", err.Error())
					continue
				}

				client, err := c.usageClientBuilder.BuildClient()
				if err != nil {
					contextutils.LoggerFrom(ctx).Errorf("Encountered error while connecting to the grpc server: %s", err.Error())
					continue
				}
				_, err = client.ReportUsage(ctx, &api.UsageRequest{
					InstanceMetadata: c.metadata,
					Payload:          payload,
				})
				if err != nil {
					contextutils.LoggerFrom(ctx).Errorf("Encountered error while reporting usage: %s", err.Error())
				}
			case <-ctx.Done():
				return
			}
		}
	}()
}
