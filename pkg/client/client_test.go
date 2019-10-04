package client

import (
	"context"
	"os"
	"time"

	"github.com/solo-io/go-utils/errors"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	v1 "github.com/solo-io/reporting-client/pkg/api/v1"
	"github.com/solo-io/reporting-client/pkg/client/mocks"
)

type testReader struct {
	payload map[string]string
}

func (t *testReader) GetPayload() (map[string]string, error) {
	return t.payload, nil
}

type testClientBuilder struct {
	client *mocks.MockReportingServiceClient
}

func (t *testClientBuilder) BuildClient() (v1.ReportingServiceClient, error) {
	return t.client, nil
}

var _ = Describe("Reporting client", func() {

	var (
		ctrl                   *gomock.Controller
		reportingServiceClient *mocks.MockReportingServiceClient
		instanceMetadata       = &v1.InstanceMetadata{
			Product: "test-product",
			Version: "v0.6.9",
			Arch:    "test-arch",
			Os:      "test-os",
		}
		testErr                = errors.New("test-err")
		ctx                    context.Context
		cancelFunc             context.CancelFunc
		pollInterval           = time.Millisecond * 50
		timeoutForEventually   = time.Second
		timeoutForConsistently = time.Millisecond * 500
	)

	var buildPayloadGetter = func(payload map[string]string) UsagePayloadReader {
		return &testReader{payload: payload}
	}

	var buildEmptyPayloadGetter = func() UsagePayloadReader {
		return buildPayloadGetter(map[string]string{})
	}

	// given a channel, build a function that forwards the request onto that channel
	var requestSender = func(reportChannel chan *v1.UsageRequest) func(ctx context.Context, request *v1.UsageRequest) (*v1.UsageResponse, error) {
		return func(ctx context.Context, request *v1.UsageRequest) (*v1.UsageResponse, error) {
			reportChannel <- request
			return &v1.UsageResponse{}, nil
		}
	}

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		reportingServiceClient = mocks.NewMockReportingServiceClient(ctrl)
		ctx, cancelFunc = context.WithCancel(context.TODO())
	})

	AfterEach(func() {
		ctrl.Finish()
		cancelFunc()
	})

	It("can report usage", func() {
		usageClient := newUsageClient(buildEmptyPayloadGetter(), instanceMetadata, func() (v1.ReportingServiceClient, error) {
			return reportingServiceClient, nil
		})

		request := &v1.UsageRequest{
			InstanceMetadata: instanceMetadata,
			Payload:          map[string]string{},
		}
		reportChannel := make(chan *v1.UsageRequest)

		reportingServiceClient.EXPECT().
			ReportUsage(ctx, request).
			DoAndReturn(requestSender(reportChannel)).
			AnyTimes()

		errorChan := usageClient.StartReportingUsage(ctx, time.Millisecond*500)

		Eventually(reportChannel, timeoutForEventually, pollInterval).Should(Receive(Equal(request)))
		Consistently(errorChan, timeoutForConsistently, pollInterval).ShouldNot(Receive())
	})

	It("reports an error on the channel if the server is unreachable", func() {
		usageClient := newUsageClient(buildEmptyPayloadGetter(), instanceMetadata, func() (v1.ReportingServiceClient, error) {
			return reportingServiceClient, nil
		})

		request := &v1.UsageRequest{
			InstanceMetadata: instanceMetadata,
			Payload:          map[string]string{},
		}
		reportChannel := make(chan *v1.UsageRequest)

		reportingServiceClient.EXPECT().
			ReportUsage(ctx, request).
			Return(nil, testErr).
			AnyTimes()

		errorChan := usageClient.StartReportingUsage(ctx, time.Millisecond*10)

		Consistently(reportChannel, timeoutForConsistently, pollInterval).ShouldNot(Receive())
		Consistently(errorChan, timeoutForConsistently, pollInterval).Should(Receive(Equal(ErrorSendingUsage(testErr))), "Should receive errors on the channel")
	})

	Context("when usage is disabled", func() {
		BeforeEach(func() {
			os.Setenv(DisableUsageVar, "true")
		})
		AfterEach(func() {
			os.Setenv(DisableUsageVar, "")
		})

		It("does not report usage", func() {
			reportChannel := make(chan *v1.UsageRequest)

			usageClient := newUsageClient(buildEmptyPayloadGetter(), instanceMetadata, func() (v1.ReportingServiceClient, error) {
				return reportingServiceClient, nil
			})

			errorChan := usageClient.StartReportingUsage(ctx, time.Millisecond*100)

			Consistently(reportChannel, timeoutForConsistently, pollInterval).ShouldNot(Receive())
			Consistently(errorChan, timeoutForConsistently, pollInterval).ShouldNot(Receive())
		})
	})
})
