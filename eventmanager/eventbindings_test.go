package eventmanager

import (
	"github.com/compozed/deployadactyl/mocks"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("EventBindings", func() {

	Describe("AddBinding", func() {
		It("should add a binding to the collection", func() {
			bindings := EventBindings{}

			bindings.AddBinding(&mocks.EventBinding{})
			bindings.AddBinding(&mocks.EventBinding{})

			Expect(len(bindings.bindings)).To(Equal(2))
		})
	})

	Describe("GetBindings", func() {
		It("should return the list of bindings", func() {
			bindings := EventBindings{}

			bindings.AddBinding(&mocks.EventBinding{})
			bindings.AddBinding(&mocks.EventBinding{})

			actual := bindings.GetBindings()

			Expect(len(actual)).To(Equal(2))
		})
	})
})
