package rate

import (
	"time"

	"github.com/redforks/testing/reset"

	bdd "github.com/onsi/ginkgo"
	"github.com/redforks/hal/timeth"
	"github.com/stretchr/testify/assert"
)

var _ = bdd.Describe("limiter", func() {
	bdd.BeforeEach(func() {
		reset.Enable()

		timeth.Install()
	})

	bdd.AfterEach(func() {
		reset.Disable()
	})

	bdd.XIt("Accept one", func() {
		l := NewLimiter(1, 10*time.Second)
		assert.True(t(), l.Accept())

		timeth.Tick(time.Second)
		assert.False(t(), l.Accept())

		timeth.Tick(9*time.Second + time.Millisecond)
		assert.True(t(), l.Accept())
	})

	bdd.It("Accept two", func() {
		l := NewLimiter(2, 10*time.Second)
		assert.True(t(), l.Accept())
		assert.True(t(), l.Accept())

		timeth.Tick(time.Second)
		assert.False(t(), l.Accept())

		timeth.Tick(9*time.Second + time.Millisecond)
		assert.True(t(), l.Accept())
		assert.True(t(), l.Accept())

		timeth.Tick(11 * time.Second)
		assert.True(t(), l.Accept())
	})

})
