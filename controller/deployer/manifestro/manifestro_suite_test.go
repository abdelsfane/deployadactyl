package manifestro_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestManifestro(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Manifestro Suite")
}
