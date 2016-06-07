package decimal_test

import (
	. "github.com/redforks/math/decimal"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("NullDecimal", func() {
	Context("Valuer", func() {
		It("null", func() {
			Ω(NullDecimal{}.Value()).Should(BeNil())
		})

		It("not null", func() {
			d := FromInt(33)
			Ω(NullDecimal{d, true}.Value()).Should(Equal("33"))
		})
	})
})
