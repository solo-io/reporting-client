package client

import (
	"context"
	"os"
	"time"

	sigmocks "github.com/solo-io/reporting-client/pkg/sig/mocks"

	"github.com/solo-io/go-utils/errors"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	v1 "github.com/solo-io/reporting-client/pkg/api/v1"
	"github.com/solo-io/reporting-client/pkg/client/mocks"
)

type testConnection struct {
}

func (t *testConnection) Close() error {
	return nil
}

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

	It("can report usage", func() {
		var (
			ctrl                   *gomock.Controller
			reportingServiceClient *mocks.MockReportingServiceClient
			product                = &v1.Product{
				Product: "test-product",
				Version: "v0.6.9",
				Arch:    "test-arch",
				Os:      "test-os",
			}
			signatureManager              *sigmocks.MockSignatureManager
			ctx                           context.Context
			cancelFunc                    context.CancelFunc
			pollInterval                  = time.Millisecond * 50
			timeoutForEventually          = time.Second
			timeoutForConsistently        = time.Millisecond * 500
			reportingServiceClientBuilder = func() (v1.ReportingServiceClient, CloseableConnection, error) {
				return reportingServiceClient, &testConnection{}, nil
			}
			signature        = "test-signature"
			errorSendTimeout = time.Millisecond * 10
		)

		ctrl = gomock.NewController(GinkgoT())
		defer ctrl.Finish()
		reportingServiceClient = mocks.NewMockReportingServiceClient(ctrl)
		signatureManager = sigmocks.NewMockSignatureManager(ctrl)
		ctx, cancelFunc = context.WithCancel(context.TODO())
		defer cancelFunc()

		signatureManager.EXPECT().
			GetSignature().
			Return(signature, nil).
			AnyTimes()

		usageClient := newUsageClient(buildEmptyPayloadGetter(), product, signatureManager, reportingServiceClientBuilder, errorSendTimeout)

		request := &v1.UsageRequest{
			InstanceMetadata: &v1.InstanceMetadata{
				Product:   product,
				Signature: signature,
			},
			Payload: map[string]string{},
		}

		// need this channel
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
		var (
			ctrl                   *gomock.Controller
			reportingServiceClient *mocks.MockReportingServiceClient
			product                = &v1.Product{
				Product: "test-product",
				Version: "v0.6.9",
				Arch:    "test-arch",
				Os:      "test-os",
			}
			signatureManager              *sigmocks.MockSignatureManager
			testErr                       = errors.New("test-err")
			ctx                           context.Context
			pollInterval                  = time.Millisecond * 50
			timeoutForConsistently        = time.Millisecond * 500
			reportingServiceClientBuilder = func() (v1.ReportingServiceClient, CloseableConnection, error) {
				return reportingServiceClient, &testConnection{}, nil
			}
			signature        = "test-signature"
			errorSendTimeout = time.Millisecond * 10
		)

		ctrl = gomock.NewController(GinkgoT())
		defer ctrl.Finish()
		reportingServiceClient = mocks.NewMockReportingServiceClient(ctrl)
		signatureManager = sigmocks.NewMockSignatureManager(ctrl)
		ctx, cancelFunc := context.WithCancel(context.TODO())
		defer cancelFunc()

		signatureManager.EXPECT().
			GetSignature().
			Return(signature, nil).
			AnyTimes()

		usageClient := newUsageClient(buildEmptyPayloadGetter(), product, signatureManager, reportingServiceClientBuilder, errorSendTimeout)

		request := &v1.UsageRequest{
			InstanceMetadata: &v1.InstanceMetadata{
				Product:   product,
				Signature: signature,
			},
			Payload: map[string]string{},
		}
		reportChannel := make(chan *v1.UsageRequest)

		reportingServiceClient.EXPECT().
			ReportUsage(ctx, request).
			Return(nil, testErr).
			AnyTimes()

		errorChan := usageClient.StartReportingUsage(ctx, time.Millisecond*10)

		Consistently(reportChannel, timeoutForConsistently, pollInterval).ShouldNot(Receive())
		Consistently(func() string {
			err := <-errorChan
			return err.Error()
		}, timeoutForConsistently, pollInterval).Should(Equal(ErrorSendingUsage(testErr).Error()), "Should receive errors on the channel")
	})

	It("still reports when the signature cannot be determined", func() {
		var (
			ctrl                   *gomock.Controller
			reportingServiceClient *mocks.MockReportingServiceClient
			product                = &v1.Product{
				Product: "test-product",
				Version: "v0.6.9",
				Arch:    "test-arch",
				Os:      "test-os",
			}
			signatureManager              *sigmocks.MockSignatureManager
			testErr                       = errors.New("test-err")
			ctx                           context.Context
			pollInterval                  = time.Millisecond * 50
			timeoutForEventually          = time.Second
			reportingServiceClientBuilder = func() (v1.ReportingServiceClient, CloseableConnection, error) {
				return reportingServiceClient, &testConnection{}, nil
			}
			signature        = "test-signature"
			errorSendTimeout = time.Millisecond * 10
		)

		ctrl = gomock.NewController(GinkgoT())
		defer ctrl.Finish()
		reportingServiceClient = mocks.NewMockReportingServiceClient(ctrl)
		signatureManager = sigmocks.NewMockSignatureManager(ctrl)
		ctx, cancelFunc := context.WithCancel(context.TODO())
		defer cancelFunc()

		signatureManager.EXPECT().
			GetSignature().
			Return(signature, nil).
			AnyTimes()

		brokenSignatureManager := sigmocks.NewMockSignatureManager(ctrl)
		brokenSignatureManager.EXPECT().
			GetSignature().
			Return("", testErr).
			MinTimes(1)

		usageClient := newUsageClient(buildEmptyPayloadGetter(), product, brokenSignatureManager, reportingServiceClientBuilder, errorSendTimeout)

		request := &v1.UsageRequest{
			InstanceMetadata: &v1.InstanceMetadata{
				Product:   product,
				Signature: "",
			},
			Payload: map[string]string{},
		}

		// need this channel
		reportChannel := make(chan *v1.UsageRequest)

		reportingServiceClient.EXPECT().
			ReportUsage(gomock.Any(), request).
			DoAndReturn(requestSender(reportChannel)).
			AnyTimes()

		errorChan := usageClient.StartReportingUsage(ctx, time.Millisecond*10)
		err := <-errorChan

		Expect(err).To(Equal(ErrorGettingSignature(testErr)))
		Eventually(reportChannel, timeoutForEventually, pollInterval).Should(Receive(Equal(request)))
	})

	Context("when usage is disabled", func() {
		BeforeEach(func() {
			os.Setenv(DisableUsageVar, "true")
		})
		AfterEach(func() {
			os.Setenv(DisableUsageVar, "")
		})

		It("does not report usage", func() {
			var (
				ctrl                   *gomock.Controller
				reportingServiceClient *mocks.MockReportingServiceClient
				product                = &v1.Product{
					Product: "test-product",
					Version: "v0.6.9",
					Arch:    "test-arch",
					Os:      "test-os",
				}
				signatureManager              *sigmocks.MockSignatureManager
				ctx                           context.Context
				pollInterval                  = time.Millisecond * 50
				reportingServiceClientBuilder = func() (v1.ReportingServiceClient, CloseableConnection, error) {
					return reportingServiceClient, &testConnection{}, nil
				}
				errorSendTimeout = time.Millisecond * 10
			)
			ctrl = gomock.NewController(GinkgoT())
			defer ctrl.Finish()
			reportingServiceClient = mocks.NewMockReportingServiceClient(ctrl)
			signatureManager = sigmocks.NewMockSignatureManager(ctrl)
			ctx, cancelFunc := context.WithCancel(context.TODO())
			defer cancelFunc()

			reportChannel := make(chan *v1.UsageRequest)

			usageClient := newUsageClient(buildEmptyPayloadGetter(), product, signatureManager, reportingServiceClientBuilder, errorSendTimeout)

			errorChan := usageClient.StartReportingUsage(ctx, time.Millisecond*100)

			Consistently(reportChannel, time.Millisecond*500, pollInterval).ShouldNot(Receive())
			Consistently(errorChan, time.Millisecond*500, pollInterval).ShouldNot(Receive())
		})
	})
})
