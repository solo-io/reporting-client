package main

import (
	"context"
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
	serverAddr := "localhost:3000"
	client, err := client.NewUsageClient(serverAddr, &testPayloadReader{}, &v1.InstanceMetadata{
		Product: "test",
		Version: "0.0.1",
		Arch:    "test",
		Os:      "test",
	})
	if err != nil {
		panic(err)
	}
	client.StartReportingUsage(context.Background(), time.Second*2)
	time.Sleep(time.Hour * 100)
}
