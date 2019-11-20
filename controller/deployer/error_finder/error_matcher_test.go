package error_finder

import (
	"github.com/compozed/deployadactyl/structs"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("RegExErrorMatcher", func() {
	It("should match no errors", func() {
		factory := ErrorMatcherFactory{}
		errorMatcher, _ := factory.CreateErrorMatcher(structs.ErrorMatcherDescriptor{
			Description: "this shouldn't bring back errors",
			Pattern:     ".{1,10}regex stuff.{1,20}",
		})
		err := errorMatcher.Match([]byte("this does not contain the regex"))
		Expect(err).To(BeNil())
	})

	It("should match one error", func() {
		factory := ErrorMatcherFactory{}
		errorMatcher, _ := factory.CreateErrorMatcher(structs.ErrorMatcherDescriptor{
			Description: "this should bring back one error",
			Pattern:     ".{1,10}ab.{1,20}",
		})
		err := errorMatcher.Match([]byte("xxxxxabxxxxxxx"))
		Expect(len(err.Details())).To(Equal(1))
		Expect(err.Error()).To(Equal("this should bring back one error"))
		Expect(err.Details()[0]).To(Equal("xxxxxabxxxxxxx"))
	})

	It("should match multiple errors", func() {
		factory := ErrorMatcherFactory{}
		errorMatcher, _ := factory.CreateErrorMatcher(structs.ErrorMatcherDescriptor{
			Description: "this should bring back multiple errors",
			Pattern:     "(?i)ab[^ab]{1,5}",
		})
		err := errorMatcher.Match([]byte("xxxxxabxxxAbxxXxxxxxxxxxxxxxxxxabx"))
		Expect(len(err.Details())).To(Equal(3))
		Expect(err.Error()).To(Equal("this should bring back multiple errors"))
		Expect(err.Details()[0]).To(Equal("abxxx"))
		Expect(err.Details()[1]).To(Equal("AbxxXxx"))
		Expect(err.Details()[2]).To(Equal("abx"))
	})

	It("should return the description", func() {
		factory := ErrorMatcherFactory{}
		errorMatcher, _ := factory.CreateErrorMatcher(structs.ErrorMatcherDescriptor{
			Description: "the description",
			Pattern:     "a regex pattern",
			Solution:    "a solution",
			Code:        "a code",
		})
		Expect(errorMatcher.Descriptor()).To(Equal("the description: a regex pattern: a solution: a code"))
	})

	It("should throw an error if pattern is missing", func() {
		factory := ErrorMatcherFactory{}
		_, err := factory.CreateErrorMatcher(structs.ErrorMatcherDescriptor{
			Description: "this should bring back the description",
			Solution:    "a solution",
		})
		Expect(err.Error()).To(Equal("error matcher requires a pattern"))
	})

	It("should return a default description if description is missing", func() {
		factory := ErrorMatcherFactory{}
		errorMatcher, _ := factory.CreateErrorMatcher(structs.ErrorMatcherDescriptor{
			Pattern:  "a regex pattern",
			Solution: "a solution",
		})
		Expect(errorMatcher.Descriptor()).To(Equal("This error does not have a description.: a regex pattern: a solution: "))
	})

	It("should return a default solution if solution is missing", func() {
		factory := ErrorMatcherFactory{}
		errorMatcher, _ := factory.CreateErrorMatcher(structs.ErrorMatcherDescriptor{
			Pattern:     "a regex pattern",
			Description: "a description",
		})
		Expect(errorMatcher.Descriptor()).To(Equal("a description: a regex pattern: No recommended solution available.: "))
	})
})
