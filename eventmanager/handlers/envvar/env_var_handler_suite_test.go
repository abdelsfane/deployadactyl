package envvar_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestEnvvarhandler(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Environment Variables Handler Suite")
}
