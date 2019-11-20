package routemapper_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestRoutemapper(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Routemapper Suite")
}
