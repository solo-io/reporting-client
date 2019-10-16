package client

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	api "github.com/solo-io/reporting-client/pkg/api/v1"
)

const (
	// set this env var to the string "true" to prevent usage from being reported
	DisableUsageVar = "DISABLE_USAGE_REPORTING"
	TestingUrl      = "localhost:8000"

	errorSendTimeout = time.Second * 10
)

// a type that knows how to load the usage payload you want to report
type UsagePayloadReader interface {
	GetPayload() (map[string]string, error)
}

type CloseableConnection interface {
	Close() error
}

// the type grpc.ClientConn is a struct- to make testing easier, hide it behind this interface
var _ CloseableConnection = &grpc.ClientConn{}

type ReportingServiceClientBuilder func() (api.ReportingServiceClient, CloseableConnection, error)

var defaultReportingServiceClientBuilder = func(url string) ReportingServiceClientBuilder {
	return func() (api.ReportingServiceClient, CloseableConnection, error) {
		var clientConn *grpc.ClientConn
		var err error

		// for testing purposes
		if url == TestingUrl {
			clientConn, err = grpc.Dial(url, grpc.WithInsecure())
		} else {
			certPool, err := x509.SystemCertPool()
			if err != nil {
				return nil, nil, err
			}
			tlsConfig := &tls.Config{
				RootCAs: certPool,
			}
			clientConn, err = grpc.Dial(url, grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)))
		}

		if err != nil {
			return nil, nil, err
		}

		client := api.NewReportingServiceClient(clientConn)
		return client, clientConn, nil
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
	errorChan          chan error
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
		errorChan:          make(chan error),
	}
}

func (c *client) StartReportingUsage(ctx context.Context, interval time.Duration) <-chan error {
	if os.Getenv(DisableUsageVar) == "true" {
		close(c.errorChan)
		return c.errorChan
	}

	// send an initial usage report immediately
	// careful not to block this goroutine
	go c.send(ctx)

	ticker := time.NewTicker(interval)
	go func() {
		for {
			select {
			case <-ticker.C:
				c.send(ctx)
			case <-ctx.Done():
				close(c.errorChan)
				return
			}
		}
	}()

	return c.errorChan
}

func (c *client) send(ctx context.Context) {
	payload, err := c.usagePayloadReader.GetPayload()
	if err != nil {
		sendWithTimeout(c.errorChan, ErrorReadingPayload(err))
		return
	}
	client, conn, err := c.usageClientBuilder()
	if err != nil {
		sendWithTimeout(c.errorChan, ErrorConnecting(err))
		return
	} else {
		defer conn.Close()
	}
	_, err = client.ReportUsage(ctx, &api.UsageRequest{
		InstanceMetadata: c.metadata,
		Payload:          payload,
	})
	if err != nil {
		sendWithTimeout(c.errorChan, ErrorSendingUsage(err))
		return
	}
}

// we don't want to block this whole goroutine if no one is listening for errors,
// so if a receiver isn't ready after the timeout, give up and continue
func sendWithTimeout(errorChan chan<- error, err error) {
	select {
	case errorChan <- err:
	case <-time.After(errorSendTimeout):
	}
}
