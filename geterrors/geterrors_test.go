package geterrors_test

import (
	. "github.com/compozed/deployadactyl/geterrors"
	"github.com/compozed/deployadactyl/randomizer"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Geterrors", func() {
	var (
		get         func(string) string
		firstKey    string
		secondKey   string
		firstValue  string
		secondValue string
	)

	BeforeEach(func() {
		firstKey = "firstKey-" + randomizer.StringRunes(10)
		secondKey = "secondKey-" + randomizer.StringRunes(10)
		firstValue = "firstValue-" + randomizer.StringRunes(10)
		secondValue = "secondValue-" + randomizer.StringRunes(10)

		get = func(key string) string {
			data := map[string]string{
				firstKey:  firstValue,
				secondKey: secondValue,
			}
			return data[key]
		}
	})

	Context("when all keys are present", func() {
		It("returns all of the values", func() {
			g := WrapFunc(get)
			Expect(g.Get(firstKey)).To(Equal(firstValue))
			Expect(g.Get(secondKey)).To(Equal(secondValue))
			Expect(g.Err("missing keys")).ToNot(HaveOccurred())
		})
	})

	Context("when a key is missing", func() {
		It("returns an error listing all of the missing keys", func() {
			g := WrapFunc(get)
			Expect(g.Get(firstKey)).To(Equal(firstValue))
			Expect(g.Get("key2")).To(Equal(""))
			Expect(g.Get(secondKey)).To(Equal(secondValue))
			Expect(g.Get("key4")).To(Equal(""))
			Expect(g.Err("missing keys")).To(MatchError("missing keys: key2, key4"))
		})
	})
})
