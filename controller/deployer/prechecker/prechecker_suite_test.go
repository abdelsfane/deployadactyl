package prechecker_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestPrechecker(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Prechecker Suite")
}
