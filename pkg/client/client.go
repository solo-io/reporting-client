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

//go:generate mockgen -destination mocks/mock_usage_payload_reader.go -package mocks github.com/solo-io/reporting-client/pkg/client UsagePayloadReader
//go:generate mockgen -destination mocks/mock_reporting_service_client.go -package mocks github.com/solo-io/reporting-client/pkg/api/v1 ReportingServiceClient

type Client interface {
	StartReportingUsage(ctx context.Context, interval time.Duration) context.CancelFunc
}

type client struct {
	usagePayloadReader UsagePayloadReader
	usageClient        api.ReportingServiceClient
	metadata           *api.InstanceMetadata
}

var _ Client = &client{}

// initializes a connection to the grpc server
// returns an error if it is unable to dial the server
func NewUsageClient(usageServerUrl string, usagePayloadReader UsagePayloadReader, instanceMetadata *api.InstanceMetadata) (*client, error) {
	clientConn, err := grpc.Dial(usageServerUrl, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	return newUsageClient(
		usagePayloadReader,
		instanceMetadata,
		api.NewReportingServiceClient(clientConn),
	)
}

// visible for testing
func newUsageClient(
	usagePayloadReader UsagePayloadReader,
	instanceMetadata *api.InstanceMetadata,
	reportingServiceClient api.ReportingServiceClient,
) (*client, error) {
	return &client{
		usagePayloadReader: usagePayloadReader,
		usageClient:        reportingServiceClient,
		metadata:           instanceMetadata,
	}, nil
}

func (c *client) StartReportingUsage(ctx context.Context, interval time.Duration) context.CancelFunc {
	cancellableCtx, cancelFunc := context.WithCancel(ctx)
	if os.Getenv(DisableUsageVar) == "true" {
		return cancelFunc
	}

	ticker := time.NewTicker(interval)
	go func() {
		for {
			select {
			case <-ticker.C:
				payload, err := c.usagePayloadReader.GetPayload()
				if err != nil {
					contextutils.LoggerFrom(ctx).Errorf("Encountered error while reading payload: %s", err.Error())
				} else {
					_, err := c.usageClient.ReportUsage(ctx, &api.UsageRequest{
						InstanceMetadata: c.metadata,
						Payload:          payload,
					})
					if err != nil {
						contextutils.LoggerFrom(ctx).Errorf("Encountered error while reporting usage: %s", err.Error())
					}
				}
				case <-cancellableCtx.Done():
					return
			}
		}
	}()

	return cancelFunc
}
