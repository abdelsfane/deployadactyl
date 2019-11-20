package state_test

import (
	C "github.com/compozed/deployadactyl/config"
	"github.com/compozed/deployadactyl/controller/deployer"
	"github.com/compozed/deployadactyl/interfaces"
	. "github.com/compozed/deployadactyl/state"
	"github.com/compozed/deployadactyl/structs"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/op/go-logging"
	"reflect"
)

var _ = Describe("Auth Resolver", func() {
	var (
		logBuffer *gbytes.Buffer
		log       interfaces.DeploymentLogger
	)

	BeforeEach(func() {
		logBuffer = gbytes.NewBuffer()
		log = interfaces.DeploymentLogger{Log: interfaces.DefaultLogger(logBuffer, logging.DEBUG, "auth_resolver_test")}
	})

	It("should write to the log", func() {
		auth := interfaces.Authorization{
			Username: "Fake_test_Username",
			Password: "Fake_test_Password",
		}

		authResolver := AuthResolver{}
		envs := structs.Environment{Authenticate: false}

		authResolver.Resolve(auth, envs, log)
		Expect(logBuffer).To(gbytes.Say("checking for basic auth"))

	})

	Context("when username and password exist", func() {
		It("should return an auth", func() {
			auth := interfaces.Authorization{
				Username: "Fake_test_Username",
				Password: "Fake_test_Password",
			}

			authResolver := AuthResolver{}
			envs := structs.Environment{Authenticate: false}

			resolveResult, err := authResolver.Resolve(auth, envs, log)

			Expect(resolveResult.Username).To(Equal("Fake_test_Username"))
			Expect(resolveResult.Password).To(Equal("Fake_test_Password"))
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Context("when authenticate is false", func() {
		Context("When username and password do not exist ", func() {
			It("should return the system account", func() {
				config := C.Config{Username: "test_username", Password: "test_password"}

				auth := interfaces.Authorization{}

				authResolver := AuthResolver{Config: config}
				envs := structs.Environment{Authenticate: false}

				resolveResult, err := authResolver.Resolve(auth, envs, log)

				Expect(resolveResult.Username).To(Equal("test_username"))
				Expect(resolveResult.Password).To(Equal("test_password"))
				Expect(err).ToNot(HaveOccurred())
			})
		})
	})

	Context("When authenticate is true", func() {
		Context("and provided username and password do not exist", func() {
			It("should return the an error", func() {
				config := C.Config{Username: "test_username", Password: "test_password"}

				auth := interfaces.Authorization{}

				authResolver := AuthResolver{Config: config}
				envs := structs.Environment{Authenticate: true}

				_, err := authResolver.Resolve(auth, envs, log)

				Expect(err).To(HaveOccurred())
				Expect(reflect.TypeOf(err)).To(Equal(reflect.TypeOf(deployer.BasicAuthError{})))
			})
		})
	})
})
