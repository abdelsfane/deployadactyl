package courier_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestCourier(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Courier Suite")
}
