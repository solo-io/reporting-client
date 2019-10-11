package main

import (
	"context"
	"log"
	"time"

	v1 "github.com/solo-io/reporting-client/pkg/api/v1"
	"github.com/solo-io/reporting-client/pkg/client"
)

type testPayloadReader struct {
}

var _ client.UsagePayloadReader = &testPayloadReader{}

func (p *testPayloadReader) GetPayload() (map[string]string, error) {
	return map[string]string{}, nil
}

// just for testing purposes
func main() {
	client := client.NewUsageClient(client.TestingUrl, &testPayloadReader{}, &v1.InstanceMetadata{
		Product: "test",
		Version: "0.0.1",
		Arch:    "test",
		Os:      "test",
	})
	errChan := client.StartReportingUsage(context.Background(), time.Second*2)
	for err := range errChan {
		log.Printf("Error: %s", err.Error())
	}
}
