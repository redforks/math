package rate

import (
	"time"

	bdd "github.com/onsi/ginkgo"
	"github.com/stretchr/testify/assert"
)

var _ = bdd.Describe("limiter", func() {
	var (
		_now time.Time

		tick = func(d time.Duration) {
			_now = _now.Add(d)
		}
	)

	bdd.BeforeEach(func() {
		_now = time.Unix(3000, 0)
		now = func() time.Time {
			return _now
		}
	})

	bdd.AfterEach(func() {
		now = time.Now
	})

	bdd.It("Accept one", func() {
		l := NewLimiter(1, 10*time.Second)
		assert.True(t(), l.Accept())

		tick(time.Second)
		assert.False(t(), l.Accept())

		tick(9*time.Second + time.Millisecond)
		assert.True(t(), l.Accept())
	})

	bdd.It("Accept two", func() {
		l := NewLimiter(2, 10*time.Second)
		assert.True(t(), l.Accept())
		assert.True(t(), l.Accept())

		tick(time.Second)
		assert.False(t(), l.Accept())

		tick(9*time.Second + time.Millisecond)
		assert.True(t(), l.Accept())
		assert.True(t(), l.Accept())

		tick(11 * time.Second)
		assert.True(t(), l.Accept())
	})

})
