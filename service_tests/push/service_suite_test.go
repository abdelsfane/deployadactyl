package push

import (
	"os"

	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"runtime"
	"testing"
)

var (
	username string
	password string
	ospath   string
)

func TestService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Service Suite")
}

var _ = BeforeSuite(func() {
	gin.SetMode(gin.TestMode)
	ospath = os.Getenv("PATH")
	var newpath string
	dir, _ := os.Getwd()
	if runtime.GOOS == "windows" {
		newpath = dir + "\\..\\..\\bin;" + ospath
	} else {
		newpath = dir + "/../../bin:" + ospath
	}
	os.Setenv("PATH", newpath)
})

var _ = AfterSuite(func() {
	os.Setenv("CF_USERNAME", username)
	os.Setenv("CF_PASSWORD", password)
	os.Setenv("PATH", ospath)
})
