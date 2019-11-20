package artifetcher_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestArtifetcher(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Artifetcher Suite")
}
