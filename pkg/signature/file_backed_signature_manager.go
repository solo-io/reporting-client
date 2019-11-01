package signature

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

const (
	signatureFileName = "signature"
	defaultDirectory  = ".soloio"
	filePermissions   = 0644
	dirPermissions    = 0755

	fileContentsTemplate = `%s

This signature is a randomly generated UUID used to de-duplicate
alerts and version information. This signature is random, it is
not based on any personally identifiable information. To create
a new signature, you can simply delete this file at any time.
See the documentation for the software using Reporting for more
information on how to disable it.
`
)

// generate a signature and persist it on disk so that we get consistent signatures across CLI invocations
// both fields are optional to provide - the signature file will be written by default to ~/.soloio/signature
type FileBackedSignatureManager struct {
	// expected to be a path to the directory where the signature file will be written
	ConfigDir string

	SignatureFileName string
}

var _ SignatureManager = &FileBackedSignatureManager{}

func (f *FileBackedSignatureManager) GetSignature() (string, error) {
	configDir := f.ConfigDir
	if configDir == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}

		configDir = path.Join(homeDir, defaultDirectory)
	}

	fileName := f.SignatureFileName
	if fileName == "" {
		fileName = signatureFileName
	}

	signatureFilePath := path.Join(configDir, fileName)

	return f.getOrGenerateSignature(signatureFilePath)
}

func (f *FileBackedSignatureManager) getOrGenerateSignature(signatureFilePath string) (string, error) {
	if _, err := os.Stat(signatureFilePath); err != nil {
		return f.writeNewSignatureFile(signatureFilePath)
	}

	fileBytes, err := ioutil.ReadFile(signatureFilePath)
	if err != nil {
		return "", err
	}

	fileContents := string(fileBytes)

	lines := strings.Split(fileContents, "\n")

	var signature string
	if len(lines) != 0 {
		signature = lines[0]
	}

	if signature == "" {
		return f.writeNewSignatureFile(signatureFilePath)
	}

	return signature, nil
}

// returns the generated signature
func (f *FileBackedSignatureManager) writeNewSignatureFile(signatureFilePath string) (string, error) {
	signature, err := f.generateSignature()
	if err != nil {
		return "", err
	}

	dir, _ := filepath.Split(signatureFilePath)
	err = os.MkdirAll(dir, dirPermissions)
	if err != nil {
		return "", err
	}

	fileContents := fmt.Sprintf(fileContentsTemplate, signature)

	return signature, ioutil.WriteFile(signatureFilePath, []byte(fileContents), filePermissions)
}

func (f *FileBackedSignatureManager) generateSignature() (string, error) {
	newUuid, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}

	return newUuid.String(), nil
}
