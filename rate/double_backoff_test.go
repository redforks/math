package rate_test

import (
	. "github.com/redforks/math/rate"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("DoubleBackoff", func() {

	It("NewDoubleBackoff", func() {
		call := func(initial, max int) func() {
			return func() {
				NewDoubleBackoff(initial, max)
			}
		}
		Ω(call(0, 1)).Should(Panic())
		Ω(call(-1, 1)).Should(Panic())
		Ω(call(1, 1)).Should(Panic())
		Ω(call(1, -1)).Should(Panic())
	})

	It("Next", func() {
		b := NewDoubleBackoff(1, 5)
		Ω(b.Next()).Should(Equal(1))
		Ω(b.Next()).Should(Equal(2))
		Ω(b.Next()).Should(Equal(4))
		Ω(b.Next()).Should(Equal(5))
		Ω(b.Next()).Should(Equal(5))
	})

	It("Reset", func() {
		b := NewDoubleBackoff(1, 5)
		Ω(b.Next()).Should(Equal(1))
		Ω(b.Next()).Should(Equal(2))

		b.Reset()
		Ω(b.Next()).Should(Equal(1))
		Ω(b.Next()).Should(Equal(2))
	})
})
