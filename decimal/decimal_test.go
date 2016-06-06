package decimal_test

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/redforks/math/decimal"
)

var _ = Describe("Decimal", func() {
	DescribeTable("FromString", func(s string) {
		d, err := decimal.FromString(s)
		Ω(err).Should(Succeed())
		Ω(d.String()).Should(Equal(s))
	},
		Entry("zero", "0"),
		Entry("Zero with scale", "0.00"),
		Entry("Very big", "123456789012345678"),
		Entry("Very small", "0.123456789"),
		Entry("Negative", "-1.30"),
	)

	DescribeTable("FromString error", func(s, errMsg string) {
		_, err := decimal.FromString(s)
		Ω(err).Should(MatchError(errMsg))
	},
		Entry("Empty string", "", `[decimal] "" not a number`),
		Entry("Not a number", "abc", `[decimal] "abc" not a number`),
		Entry("Like a number", "1.3.3", `[decimal] "1.3.3" not a number`),
		Entry("Effective number too large", "12345678901234567890", `[decimal] "12345678901234567890" effective number out of range`),
		Entry("Scale out of range", "0.1234567890", `[decimal] scale 10 out of range`),
	)

	DescribeTable("FromStringWithScale error", func(s string, scale int, errMsg string) {
		_, err := decimal.FromStringWithScale(s, scale)
		Ω(err).Should(MatchError(errMsg))
	},
		Entry("Empty string", "", 0, `[decimal] "" not a number`),
		Entry("Not a number", "abc", 0, `[decimal] "abc" not a number`),
		Entry("Like a number", "1.2.3", 0, `[decimal] "1.2.3" not a number`),
		Entry("Effective number too large", "12345678901234567890", 0, `[decimal] "12345678901234567890" effective number out of range`),
		Entry("Actual scale out of range", "0.1234567890", 1, `[decimal] scale 10 out of range`),
		Entry("Scale out of range 1", "0.0", 10, `[decimal] scale 10 out of range`),
		Entry("Scale out of range 2", "0.0", -1, `[decimal] scale -1 out of range`),
	)

	DescribeTable("FromStringWithScale", func(s string, scale int, exp string) {
		d, err := decimal.FromStringWithScale(s, scale)
		Ω(err).Should(Succeed())
		Ω(d.String()).Should(Equal(exp))
	},
		Entry("Scale matches 1", "0", 0, "0"),
		Entry("Scale matches 2", "0.00", 2, "0.00"),
		Entry("Scale larger", "3.3", 3, "3.300"),
		Entry("Scale smaller", "3.333", 1, "3.333"),
	)

	It("GoStringer", func() {
		d, err := decimal.FromString("3.30")
		Ω(err).Should(Succeed())
		Ω(d.GoString()).Should(Equal("3.30m"))
	})

	DescribeTable("ShortString", func(str, exp string) {
		d, err := decimal.FromString(str)
		Ω(err).Should(Succeed())
		Ω(d.ShortString()).Should(Equal(exp))
	},
		Entry("No ending zero", "3.4", "3.4"),
		Entry("Ending integer zeros", "300", "300"),
		Entry("Ending fragment zeros", "3.00", "3"),
	)

	DescribeTable("FromInt", func(i int64) {
		d := decimal.FromInt(i)
		Ω(d.Scale()).Should(Equal(uint8(0)))
		Ω(d.String()).Should(Equal(fmt.Sprint(i)))
	},
		Entry("zero", int64(0)),
		Entry("negative 1", int64(-1)),
		Entry("very large", int64(123456789012345678)),
		Entry("very large negative", int64(-123456789012345678)),
	)

	DescribeTable("ToInt", func(s string, exp int64) {
		d, err := decimal.FromString(s)
		Ω(err).Should(Succeed())

		Ω(d.ToInt()).Should(Equal(exp))
	},
		Entry("No fragment", "345", int64(345)),
		Entry("No fragment negative", "-345", int64(-345)),
		Entry("large value", "123456789012345678", int64(123456789012345678)),
		Entry("With fragment", "3.345", int64(3)),
		Entry("With fragment negative", "-3.345", int64(-3)),
		Entry("Round up", "3.5", int64(4)),
		Entry("Round up negative", "-3.5", int64(-4)),
	)

	Context("Compute", func() {

		DescribeTable("Negate", func(s, exp string) {
			d, err := decimal.FromString(s)
			Ω(err).Should(Succeed())
			Ω(d.Negate().String()).Should(Equal(exp))
		},
			Entry("Zero", "0", "0"),
			Entry("Zero 2", "0.00", "0.00"),
			Entry("Positive", "3.456", "-3.456"),
			Entry("Negative", "-3.4456", "3.4456"),
		)

		DescribeTable("Add", func(a, b, c string) {
			x, err := decimal.FromString(a)
			Ω(err).Should(Succeed())
			y, err := decimal.FromString(b)
			Ω(err).Should(Succeed())

			z := x.Add(y)
			Ω(z.String()).Should(Equal(c))
		},
			Entry("Two positive integers", "3", "4", "7"),
			Entry("Positive negative integers", "3", "-4", "-1"),
			Entry("Scale equaled", "1.1", "2.2", "3.3"),
			Entry("Scale not equal", "1", "4.5", "5.5"),
			Entry("Add up to integer", "1.2", "2.8", "4.0"),
			Entry("fragment", "0.0003", "0.0007", "0.0010"),
			Entry("fragment and integer", "0.3", "300", "300.3"),
		)

		DescribeTable("AddToScale", func(a, b, c string, scale int) {
			x, err := decimal.FromString(a)
			Ω(err).Should(Succeed())
			y, err := decimal.FromString(b)
			Ω(err).Should(Succeed())

			z := x.AddToScale(y, scale)
			Ω(z.String()).Should(Equal(c))
		},
			Entry("Scale not changed", "3.10", "4.01", "7.11", 2),
			Entry("Extend scale", "3.10", "1", "4.100", 3),
			Entry("Shrink scale round up", "3.1", "0.05", "3.2", 1),
			Entry("Shrink scale round down", "3.1", "0.04", "3.1", 1),
		)

	})

	XContext("Subtract", func() {

	})

	XContext("Multiply", func() {})

	XContext("Div", func() {

	})

	DescribeTable("Round", func(s string, scale int, exp string) {
		d, err := decimal.FromString(s)
		Ω(err).Should(Succeed())
		Ω(d.Round(scale).String()).Should(Equal(exp))
	},
		Entry("Scale not changed", "3.4", 1, "3.4"),
		Entry("Expand scale", "3.4", 3, "3.400"),
		Entry("Shrink scale round up", "3.45", 1, "3.5"),
		Entry("Shrink scale round down", "3.44", 1, "3.4"),
		Entry("Shrink scale round up negative", "-3.45", 1, "-3.5"),
		Entry("Shrink scale round down negative", "-3.44", 1, "-3.4"),
	)

	XContext("GetZero", func() {
		// get all scale of zeros
	})
})
