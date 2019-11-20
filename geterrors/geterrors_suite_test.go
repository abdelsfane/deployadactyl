package geterrors_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestGeterrors(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Geterrors Suite")
}
