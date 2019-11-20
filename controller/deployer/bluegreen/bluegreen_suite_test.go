package bluegreen_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestBluegreen(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Bluegreen Suite")
}
