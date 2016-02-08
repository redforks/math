package rate

import (
	"log"
	"sync"
	"time"
)

const (
	tag = "math-rate"
)

var (
	// return current time, set to test stub during unit test.
	now func() time.Time = time.Now
)

// Limiter is a rate control that not allow things to happen more than n times
// in a specified duration.
type Limiter struct {
	l sync.Mutex

	ring []time.Time // ring buffer of last n thing happens time.
	tail int         // index of oldest item in ring

	d time.Duration
}

// NewLimiter create a limiter limits things happens no more than n times in m
// duration.
func NewLimiter(n int, m time.Duration) *Limiter {
	if n <= 0 {
		log.Panicf("[%s] n (%d) of Limiter can not less than 0", tag, n)
	}
	return &Limiter{ring: make([]time.Time, n), d: m}
}

// Accept returns true if currently can accept request.
func (l *Limiter) Accept() bool {
	l.l.Lock()

	t := now()
	r := t.Sub(l.ring[l.tail]) > l.d

	if r {
		l.ring[l.tail] = t

		l.tail--
		if l.tail == -1 {
			l.tail = len(l.ring) - 1
		}
	}

	l.l.Unlock()
	return r
}
