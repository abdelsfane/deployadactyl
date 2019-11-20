package manifestro_test

import (
	. "github.com/compozed/deployadactyl/controller/deployer/manifestro"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Manifestro", func() {

	Context("when manifest is empty", func() {
		It("returns nil", func() {
			manifest := ""

			result := GetInstances(manifest)

			Expect(result).To(BeNil())
		})
	})

	Context("when manifest not valid", func() {
		It("returns nil", func() {
			manifest := "bork"

			result := GetInstances(manifest)

			Expect(result).To(BeNil())
		})
	})

	Context("when manifest is not empty", func() {
		Context("when instances is not found", func() {
			It("returns nil", func() {
				manifest := `
applications:
- name: example`

				result := GetInstances(manifest)

				Expect(result).To(BeNil())
			})
		})

		Context("when instances is found", func() {
			It("returns the number of instances", func() {
				manifest := `
applications:
- name: example
  instances: 2`

				result := GetInstances(manifest)

				Expect(*result).To(Equal(uint16(2)))
			})
		})

		Context("when instances is zero", func() {
			It("returns nil", func() {
				manifest := `
applications:
- name: example
  instances: 0`

				result := GetInstances(manifest)

				Expect(result).To(BeNil())
			})
		})

		Context("when instances is not a number", func() {
			It("returns nil", func() {
				manifest := `
applications:
- name: example
  instances: bork`

				result := GetInstances(manifest)

				Expect(result).To(BeNil())
			})
		})

		Context("when applications is not found", func() {
			It("returns nil", func() {
				manifest := `---
env:
  COOL_DINOSAUR: majungasaurus
`
				result := GetInstances(manifest)

				Expect(result).To(BeNil())
			})
		})
	})

	Context("when instances is found", func() {
		Context("when there are multiple applications", func() {
			It("returns the number of instances from the first application", func() {
				manifest := `
applications:
- name: example
  instances: 2
- name: example2
  instances: 4`

				result := GetInstances(manifest)

				Expect(*result).To(Equal(uint16(2)))
			})
		})
	})
})
