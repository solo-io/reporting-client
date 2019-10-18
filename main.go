package main

import (
	"context"
	"log"
	"time"

	"github.com/solo-io/reporting-client/pkg/signature"

	v1 "github.com/solo-io/reporting-client/pkg/api/v1"
	"github.com/solo-io/reporting-client/pkg/client"
)

type testPayloadReader struct {
}

var _ client.UsagePayloadReader = &testPayloadReader{}

func (p *testPayloadReader) GetPayload() (map[string]string, error) {
	return map[string]string{}, nil
}

type testSignatureManager struct {
}

func (t *testSignatureManager) GetSignature() (string, error) {
	return "test-signature", nil
}

var _ signature.SignatureManager = &testSignatureManager{}

// just for testing purposes
func main() {
	client := client.NewUsageClient(
		client.TestingUrl,
		&testPayloadReader{},
		&v1.Product{
			Product: "test-product",
			Version: "0.6.9",
			Arch:    "test-arch",
			Os:      "test-os",
		},
		&testSignatureManager{},
	)
	errChan := client.StartReportingUsage(context.Background(), time.Second*2)
	for err := range errChan {
		log.Printf("Error: %s", err.Error())
	}
}
