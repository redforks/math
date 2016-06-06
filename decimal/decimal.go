package decimal

import (
	"fmt"
	"strings"

	"strconv"
)

const tag = "decimal"

// Decimal is a decimal implementation provides 18 max effective numbers.
// Decimal is immutable after creation, passed by value. Do not change its fields.
type Decimal struct {
	digits int64
	scale  uint8
}

// FromInt create Decimal from 64bit integer, set to scale to zero.
func FromInt(i int64) Decimal {
	return Decimal{i, 0}
}

// FromString create decimal from string, scale set from fragment part of number.
// Such as '3.00', scale is 2.
func FromString(s string) (Decimal, error) {
	return FromStringWithScale(s, 0)
}

// FromStringWithScale create decimal from string, with specific scale.
// Use number's actual scale if it larger than specific scale. Examples:
//
//  FromStringWithScale("3", 2) // 3.00
//  FromStringWithScale("3.33", 0) // 3.33
func FromStringWithScale(str string, scale int) (Decimal, error) {
	if err := checkScale(scale); err != nil {
		return Decimal{}, err
	}

	s := str
	actScale, dotIdx := 0, strings.IndexRune(s, '.')
	if dotIdx != -1 {
		actScale = len(s) - 1 - dotIdx
		s = s[0:dotIdx] + s[dotIdx+1:]
	}

	digits, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		e := err.(*strconv.NumError)
		switch e.Err {
		case strconv.ErrSyntax:
			return Decimal{}, fmt.Errorf("[%s] \"%s\" not a number", tag, str)
		case strconv.ErrRange:
			return Decimal{}, fmt.Errorf("[%s] \"%s\" effective number out of range", tag, str)
		default:
			return Decimal{}, err
		}
	}

	if scale > actScale {
		digits *= powerOf10(scale - actScale)
	} else {
		scale = actScale
	}
	if err := checkScale(scale); err != nil {
		return Decimal{}, err
	}

	// TODO: scale out of range > 10 < 0
	return Decimal{digits, uint8(scale)}, err
}

// String implement fmt.Stringer interface, return decimal value in string format,
// appended 0 to scales. Such as 3.00, use ShortString() to get '3'.
func (d Decimal) String() string {
	if d.scale == 0 {
		return strconv.FormatInt(d.digits, 10)
	}

	var s string
	isNegative := d.digits < 0
	if isNegative {
		s = strconv.FormatInt(-d.digits, 10)
	} else {
		s = strconv.FormatInt(d.digits, 10)
	}

	dotIdx := len(s) - int(d.scale)
	switch {
	case dotIdx == 0:
		s = "0." + s
	case dotIdx < 0:
		s = "0." + strings.Repeat("0", -dotIdx) + s
	default:
		s = s[0:dotIdx] + "." + s[dotIdx:]
	}

	if isNegative {
		s = "-" + s
	}
	return s
}

// GoString implement fmt.GoStringer interface. Adding 'm' suffix to result of String().
func (d Decimal) GoString() string {
	return d.String() + "m"
}

// ShortString convert current value to string, removing ending 0s.
// Such as 3.00, returns 3.
func (d Decimal) ShortString() string {
	r := d.String()
	if d.scale == 0 {
		return r
	}

	return strings.TrimRight(r, ".0")
}

// Scale return scale of this decimal value.
func (d Decimal) Scale() uint8 {
	return d.scale
}

// ToInt convert current value to int, round tenth fragment.
func (d Decimal) ToInt() int64 {
	if d.scale == 0 {
		return d.digits
	}

	r := d.digits
	r /= powerOf10(int(d.scale - 1))
	return roundLastDecimalBit(r)
}

// Round decimal to specific scale.
func (d Decimal) Round(scale int) Decimal {
	if err := checkScale(scale); err != nil {
		panic(err.Error())
	}

	diff := scale - int(d.scale)
	digits := d.digits
	switch {
	case diff == 0:
		return d
	case diff > 0:
		digits *= powerOf10(diff)
	default:
		diff = -diff
		digits /= powerOf10(diff - 1)
		digits = roundLastDecimalBit(digits)
	}

	return Decimal{digits, uint8(scale)}
}

// Negate returns negative value
func (d Decimal) Negate() Decimal {
	return Decimal{-d.digits, d.scale}
}

// Add this value with other value, use two values' highest scale as result scale, such as
// 3.45 + 1 = 4.45.
func (d Decimal) Add(other Decimal) Decimal {
	va, vb, scale := d.digits, other.digits, d.scale
	diff := int(d.scale) - int(other.scale)
	switch {
	case diff > 0:
		vb = vb * powerOf10(diff)
	case diff < 0:
		scale = other.scale
		va = va * powerOf10(-diff)
	}
	return Decimal{va + vb, scale}
}

// AddToScale this value with other value round to specific scale.
func (d Decimal) AddToScale(other Decimal, scale int) Decimal {
	return d.Add(other).Round(scale)
}

var tenth = [9]int64{1, 10, 100, 1000, 10000, 100000, 1000000, 10000000, 100000000}

func powerOf10(n int) int64 {
	return tenth[n]
}

// checkScale checks scale, return non-nil error if out of range
func checkScale(scale int) error {
	if scale >= 10 || scale < 0 {
		return fmt.Errorf("[%s] scale %d out of range", tag, scale)
	}
	return nil
}

// roundLastDecimalBit round last decimal bit and divide by 10.
func roundLastDecimalBit(v int64) int64 {
	mod, v := v%10, v/10
	switch {
	case mod >= 5:
		v++
	case mod <= -5:
		v--
	}
	return v
}
