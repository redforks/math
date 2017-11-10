package rate

import "fmt"

// DoubleBackOff implement a Backoff logic, that given a start value, then next
// value is twice of previous value, until max value reached..
type DoubleBackoff struct {
	initial, max, cur int

	// backoff times
	times int
}

// NewDoubleBackoff create a new instance of DoubleBackoff, panic if:
// initial value less or equals to zero, max value less or equals initial.
func NewDoubleBackoff(initial, max int) *DoubleBackoff {
	if initial <= 0 {
		panic(fmt.Sprintf("[%s] invalid DoubleBackoff(%d, %d) initial value", tag, initial, max))
	}

	if max <= initial {
		panic(fmt.Sprintf("[%s] invalid DoubleBackoff(%d, %d) max value", tag, initial, max))
	}
	return &DoubleBackoff{
		initial: initial,
		max:     max,
		cur:     initial,
	}
}

func (b *DoubleBackoff) Next() int {
	if b.times != 0 {
		b.cur *= 2
	}
	b.times++
	if b.cur > b.max {
		b.cur = b.max
	}
	return b.cur
}

func (b *DoubleBackoff) Reset() {
	b.times = 0
	b.cur = b.initial
}
