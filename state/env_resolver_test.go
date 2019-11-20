package state_test

import (
	DC "github.com/compozed/deployadactyl/config"
	"github.com/compozed/deployadactyl/controller/deployer"
	. "github.com/compozed/deployadactyl/state"
	"github.com/compozed/deployadactyl/structs"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"reflect"
)

var _ = Describe("Env Resolver", func() {

	Context("when Environment Exist", func() {
		It("should return an environment object", func() {
			testEnvString := "Env-Fake-Data"

			config := DC.Config{Environments: make(map[string]structs.Environment, 0)}
			config.Environments["Env-Fake-Data"] = structs.Environment{
				Name: "Env-Fake-Data",
			}
			envResolver := EnvResolver{Config: config}

			resolveResult, err := envResolver.Resolve(testEnvString)

			Expect(resolveResult).To(Equal(config.Environments["Env-Fake-Data"]))
			Expect(resolveResult.Name).To(Equal("Env-Fake-Data"))
			Expect(err).ToNot(HaveOccurred())
		})
	})
	Context("When Environment does not exist", func() {
		It("should return an error", func() {
			testString := "Env-Fake-Data"

			config := DC.Config{Environments: make(map[string]structs.Environment, 0)}

			envResolver := EnvResolver{Config: config}
			_, err := envResolver.Resolve(testString)
			Expect(err).To(HaveOccurred())
			Expect(reflect.TypeOf(err)).To(Equal(reflect.TypeOf(deployer.EnvironmentNotFoundError{})))

		})
	})

})
