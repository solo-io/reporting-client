package signature

import (
	"io/ioutil"
	"os"
	"path/filepath"
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

	Context("In-memory signature manager", func() {
		It("can generate a new signature when first created", func() {
			sigManager := NewSignatureManager()
			signature, err := sigManager.GetSignature()
			Expect(err).NotTo(HaveOccurred())

			Expect(signature).NotTo(BeEmpty())
		})

		It("can regenerate a signature if the original is unrecoverable", func() {
			sigManager := inMemorySignatureManager{
				mutex: sync.Mutex{},
			}

			signature, err := sigManager.GetSignature()
			Expect(err).NotTo(HaveOccurred())
			Expect(signature).NotTo(BeEmpty())

			original := signature

			sigManager.signature = ""

			newSig, err := sigManager.GetSignature()
			Expect(err).NotTo(HaveOccurred())
			Expect(newSig).NotTo(BeEmpty())

			Expect(original).NotTo(Equal(newSig))
		})
	})

	Context("File-backed signature manager", func() {
		It("can create a signature", func() {
			tempDir := os.TempDir()

			tempFilePath, err := ioutil.TempFile(tempDir, "")
			Expect(err).NotTo(HaveOccurred(), "Should be able to get a temp file")
			fileName, _ := filepath.Split(tempFilePath.Name())

			manager := &FileBackedSignatureManager{
				ConfigDir:         tempDir,
				SignatureFileName: fileName,
			}

			signature, err := manager.GetSignature()
			Expect(err).NotTo(HaveOccurred(), "Should be able to get a signature")
			Expect(signature).NotTo(BeEmpty(), "The signature should not be empty")
			Expect(os.Remove(tempFilePath.Name())).NotTo(HaveOccurred())
		})

		It("can regenerate a signature", func() {
			tempDir := os.TempDir()

			tempFilePath, err := ioutil.TempFile(tempDir, "")
			Expect(err).NotTo(HaveOccurred(), "Should be able to get a temp file")
			_, fileName := filepath.Split(tempFilePath.Name())

			manager := &FileBackedSignatureManager{
				ConfigDir:         tempDir,
				SignatureFileName: fileName,
			}

			signature, err := manager.GetSignature()
			Expect(err).NotTo(HaveOccurred(), "Should be able to get a signature")
			Expect(signature).NotTo(BeEmpty(), "The signature should not be empty")

			Expect(os.Remove(tempFilePath.Name())).NotTo(HaveOccurred(), "Should remove the first temp file")

			tempFilePath, err = ioutil.TempFile(tempDir, "")
			Expect(err).NotTo(HaveOccurred(), "Should be able to get a NEW temp file")
			_, fileName = filepath.Split(tempFilePath.Name())

			manager.SignatureFileName = fileName

			newSignature, err := manager.GetSignature()
			Expect(err).NotTo(HaveOccurred(), "Should be able to generate the second signature")
			Expect(newSignature).NotTo(Equal(signature), "Should get a different signature than we had before")
			Expect(newSignature).NotTo(BeEmpty(), "Should generate a nonempty signature")

			Expect(os.Remove(tempFilePath.Name())).NotTo(HaveOccurred(), "Should remove the second temp file")
		})
	})
})
