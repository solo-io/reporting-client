package client

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"os"
	"sync"
	"time"

	"github.com/solo-io/reporting-client/pkg/sig"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	api "github.com/solo-io/reporting-client/pkg/api/v1"
)

const (
	// set this env var to the string "true" to prevent usage from being reported
	DisableUsageVar = "DISABLE_USAGE_REPORTING"
	TestingUrl      = "localhost:3000"

	defaultErrorSendTimeout = time.Second * 10
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
	product            *api.Product
	signatureManager   sig.SignatureManager
	errorChan          chan error
	errorSendTimeout   time.Duration

	// we have goroutines and timeouts and etc so things are messy- lock before
	// interacting with the error channel. Otherwise in testing you'll see data races happening
	errorChannelMutex sync.Mutex
}

var _ Client = &client{}

func NewUsageClient(usageServerUrl string, usagePayloadReader UsagePayloadReader, product *api.Product, signatureManager sig.SignatureManager) *client {
	return newUsageClient(
		usagePayloadReader,
		product,
		signatureManager,
		defaultReportingServiceClientBuilder(usageServerUrl),
		defaultErrorSendTimeout,
	)
}

// visible for testing
func newUsageClient(
	usagePayloadReader UsagePayloadReader,
	product *api.Product,
	signatureManager sig.SignatureManager,
	reportingServiceClientBuilder ReportingServiceClientBuilder,
	errorSendTimeout time.Duration,
) *client {
	return &client{
		usagePayloadReader: usagePayloadReader,
		usageClientBuilder: reportingServiceClientBuilder,
		product:            product,
		signatureManager:   signatureManager,
		errorChan:          make(chan error),
		errorSendTimeout:   errorSendTimeout,
		errorChannelMutex:  sync.Mutex{},
	}
}

func (c *client) StartReportingUsage(ctx context.Context, interval time.Duration) <-chan error {
	if os.Getenv(DisableUsageVar) == "true" {
		close(c.errorChan)
		return c.errorChan
	}

	// send an initial usage report immediately
	// careful not to block this goroutine
	go c.sendUsage(ctx)

	ticker := time.NewTicker(interval)
	go func() {
		for {
			select {
			case <-ticker.C:
				c.sendUsage(ctx)
			case <-ctx.Done():
				c.errorChannelMutex.Lock()
				close(c.errorChan)
				c.errorChannelMutex.Unlock()
				return
			}
		}
	}()

	return c.errorChan
}

func (c *client) sendUsage(ctx context.Context) {
	payload, err := c.usagePayloadReader.GetPayload()
	if err != nil {
		c.sendErrorWithTimeout(ctx, ErrorReadingPayload(err))
		return
	}
	client, conn, err := c.usageClientBuilder()
	if err != nil {
		c.sendErrorWithTimeout(ctx, ErrorConnecting(err))
		return
	} else {
		defer conn.Close()
	}
	signature, err := c.signatureManager.GetSignature()
	if err != nil {
		// we still want to report usage even if the signature is busted, so don't return early here
		c.sendErrorWithTimeout(ctx, ErrorGettingSignature(err))
	}

	_, err = client.ReportUsage(ctx, &api.UsageRequest{
		InstanceMetadata: &api.InstanceMetadata{
			Product:   c.product,
			Signature: signature,
		},
		Payload: payload,
	})
	if err != nil {
		c.sendErrorWithTimeout(ctx, ErrorSendingUsage(err))
		return
	}
}

// we don't want to block this whole goroutine if no one is listening for errors,
// so if a receiver isn't ready after the timeout, give up and continue
func (c *client) sendErrorWithTimeout(ctx context.Context, err error) {
	c.errorChannelMutex.Lock()
	defer c.errorChannelMutex.Unlock()

	select {
	case <-ctx.Done():
	case c.errorChan <- err:
	case <-time.After(c.errorSendTimeout):
	}
}
