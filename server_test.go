package main_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

var goodConfig = []byte(`---
environments:
  - name: test
    domain: examples.are.cool.com
    authenticate: false
    skip_ssl: true
    instances: 2
    foundations:
    - https://example.endpoint1.cf.com
    - https://example.endpoint2.cf.com
`)

var badConfig = []byte(`---
environments:
  - name: sandbox
`)

var _ = Describe("Server", func() {

	var (
		session *gexec.Session
		err     error
	)
	BeforeEach(func() {
		os.Setenv("CF_USERNAME", "test user")
		os.Setenv("CF_PASSWORD", "test pwd")
	})

	AfterEach(func() {
		os.Unsetenv("CF_USERNAME")
		os.Unsetenv("CF_PASSWORD")
		session.Terminate()
	})

	Describe("log level flag", func() {
		Context("when a log level is not specified", func() {
			It("uses the default log level ", func() {
				level := os.Getenv("DEPLOYADACTYL_LOGLEVEL")

				os.Unsetenv("DEPLOYADACTYL_LOGLEVEL")
				Expect(err).ToNot(HaveOccurred())

				session, err = gexec.Start(exec.Command(pathToCLI), GinkgoWriter, GinkgoWriter)
				Expect(err).ToNot(HaveOccurred())

				Eventually(session.Out).Should(Say("log level"))
				Eventually(session.Out).Should(Say("DEBUG"))

				os.Setenv("DEPLOYADACTYL_LOGLEVEL", level)
			})
		})

		Context("when log level is invalid", func() {
			It("throws an error", func() {
				level := os.Getenv("DEPLOYADACTYL_LOGLEVEL")

				Expect(os.Setenv("DEPLOYADACTYL_LOGLEVEL", "tanystropheus")).To(Succeed())

				session, err = gexec.Start(exec.Command(pathToCLI), GinkgoWriter, GinkgoWriter)
				Expect(err).ToNot(HaveOccurred())

				Eventually(session.Err).Should(Say("invalid log level"))

				os.Setenv("DEPLOYADACTYL_LOGLEVEL", level)
			})
		})
	})

	Describe("command line flags", func() {
		Describe("config flag", func() {
			Context("when the config flag is not provided", func() {
				It("throws an error", func() {
					session, err = gexec.Start(exec.Command(pathToCLI), GinkgoWriter, GinkgoWriter)
					Expect(err).ToNot(HaveOccurred())

					Eventually(session.Out).Should(Say("no such file or directory"))
				})
			})

			Context("when an invalid config path is specified", func() {
				It("throws an error", func() {
					session, err = gexec.Start(exec.Command(pathToCLI, "-config", "./gorgosaurus.yml"), GinkgoWriter, GinkgoWriter)
					Expect(err).ToNot(HaveOccurred())

					Eventually(session.Out).Should(Say("no such file or directory"))
				})
			})

			Context("when a bad config is provided", func() {
				It("returns an error", func() {
					configLocation := fmt.Sprintf("%s/config.yml", path.Dir(pathToCLI))

					Expect(ioutil.WriteFile(configLocation, badConfig, 0777)).To(Succeed())

					session, err = gexec.Start(exec.Command(pathToCLI, "-config", configLocation), GinkgoWriter, GinkgoWriter)
					Expect(err).ToNot(HaveOccurred())

					Eventually(session.Out).Should(Say("missing required parameter"))
				})
			})
		})
	})
})
