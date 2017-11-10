package rate_test

import (
	"time"

	"github.com/redforks/testing/reset"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/redforks/hal/timeth"
	. "github.com/redforks/math/rate"
)

var _ = Describe("limiter", func() {
	BeforeEach(func() {
		reset.Enable()

		timeth.Install()
	})

	AfterEach(func() {
		reset.Disable()
	})

	It("Accept one", func() {
		l := NewLimiter(1, 10*time.Second)
		Ω(l.Accept()).Should(BeTrue())

		timeth.Tick(time.Second)
		Ω(l.Accept()).Should(BeFalse())

		timeth.Tick(9*time.Second + time.Millisecond)
		Ω(l.Accept()).Should(BeTrue())
	})

	It("Accept two", func() {
		l := NewLimiter(2, 10*time.Second)
		Ω(l.Accept()).Should(BeTrue())
		Ω(l.Accept()).Should(BeTrue())

		timeth.Tick(time.Second)
		Ω(l.Accept()).Should(BeFalse())

		timeth.Tick(9*time.Second + time.Millisecond)
		Ω(l.Accept()).Should(BeTrue())
		Ω(l.Accept()).Should(BeTrue())

		timeth.Tick(11 * time.Second)
		Ω(l.Accept()).Should(BeTrue())
	})

})
