package sig

import (
	"sync"

	"github.com/google/uuid"
)

const ErrorSignature = "ERROR"

// users of reporting-client are encouraged but not required to use this type to keep track of their signature
type SignatureManager interface {

	// get the signature for this reporting client instance
	// this function may return both a nonempty string and an error in the case
	// where we failed to load the previously-existing signature but generated a new one.
	// It may return the const ErrorSignature when it is totally unable to generate a signature
	GetSignature() (string, error)
}

func NewSignatureManager() SignatureManager {
	return &inMemorySignatureManager{
		mutex: sync.Mutex{},
	}
}

type inMemorySignatureManager struct {
	signature string
	mutex     sync.Mutex
}

func (i *inMemorySignatureManager) GetSignature() (string, error) {
	i.optionallyRegenerateSignature()

	return i.signature, nil
}

func (i *inMemorySignatureManager) optionallyRegenerateSignature() {
	i.mutex.Lock()
	defer i.mutex.Unlock()

	if i.signature == "" {
		signature, err := uuid.NewRandom()

		if err == nil {
			i.signature = signature.String()
		} else {
			i.signature = ErrorSignature
		}
	}
}