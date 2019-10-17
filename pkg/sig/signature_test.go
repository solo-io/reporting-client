package sig

import (
	"sync"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Signature Manager", func() {

	var (
		ctrl *gomock.Controller
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	It("can generate a new signature when first created", func() {
		sigManager := NewSignatureManager()
		signature, err := sigManager.GetSignature()
		Expect(err).NotTo(HaveOccurred())

		Expect(signature).NotTo(BeEmpty())
		Expect(signature).NotTo(Equal(ErrorSignature))
	})

	It("can regenerate a signature if the original is unrecoverable", func() {
		sigManager := inMemorySignatureManager{
			mutex: sync.Mutex{},
		}

		signature, err := sigManager.GetSignature()
		Expect(err).NotTo(HaveOccurred())
		Expect(signature).NotTo(BeEmpty())
		Expect(signature).NotTo(Equal(ErrorSignature))

		original := signature

		sigManager.signature = ""

		newSig, err := sigManager.GetSignature()
		Expect(err).NotTo(HaveOccurred())
		Expect(newSig).NotTo(BeEmpty())
		Expect(newSig).NotTo(Equal(ErrorSignature))

		Expect(original).NotTo(Equal(newSig))
	})
})
