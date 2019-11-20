package error_finder_test

import (
	. "github.com/compozed/deployadactyl/controller/deployer/error_finder"

	"github.com/compozed/deployadactyl/interfaces"
	"github.com/compozed/deployadactyl/mocks"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ErrorFinder", func() {

	It("returns no errors when no matchers are configured", func() {
		errorFinder := ErrorFinder{}
		errors := errorFinder.FindErrors("This is some text that doesn't affect the test")

		Expect(len(errors)).To(BeZero())
	})

	It("returns multiple errors when matchers are configured", func() {
		matchers := make([]interfaces.ErrorMatcher, 0, 0)

		matcher := &mocks.ErrorMatcherMock{}
		matcher.MatchCall.Returns = CreateLogMatchedError("a test error", []string{"error 1", "error 2", "error 3"}, "error solution", "test code")
		matchers = append(matchers, matcher)

		matcher = &mocks.ErrorMatcherMock{}
		matcher.MatchCall.Returns = CreateLogMatchedError("another test error", []string{"error 4", "error 5", "error 6"}, "another error solution", "another test code")
		matchers = append(matchers, matcher)

		errorFinder := ErrorFinder{Matchers: matchers}
		errors := errorFinder.FindErrors("This is some text that doesn't affect the test")

		Expect(len(errors)).To(Equal(2))
		Expect(errors[0].Error()).To(Equal("a test error"))
		Expect(errors[0].Details()[0]).To(Equal("error 1"))
		Expect(errors[0].Solution()).To(Equal("error solution"))
		Expect(errors[0].Code()).To(Equal("test code"))
		Expect(errors[1].Error()).To(Equal("another test error"))
		Expect(errors[1].Details()[2]).To(Equal("error 6"))
		Expect(errors[1].Solution()).To(Equal("another error solution"))
		Expect(errors[1].Code()).To(Equal("another test code"))
	})

})
