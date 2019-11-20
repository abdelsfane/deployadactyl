package creator

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestCreator(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Creator Suite")
}
