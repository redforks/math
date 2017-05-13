package decimal_test

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/redforks/math/decimal"
	"github.com/redforks/testing/matcher"
)

var _ = Describe("Decimal", func() {
	toDecimal2 := func(a, b string) (x, y decimal.Decimal) {
		Ω(decimal.FromString(a)).Should(matcher.Save(&x))
		Ω(decimal.FromString(b)).Should(matcher.Save(&y))
		return
	}

	assertBinOp := func(a, b string, expected interface{}, op func(x, y decimal.Decimal) interface{}) {
		if s, ok := expected.(string); ok {
			var d decimal.Decimal
			Ω(decimal.FromString(s)).Should(matcher.Save(&d))
			expected = d
		}
		x, y := toDecimal2(a, b)
		Ω(op(x, y)).Should(Equal(expected))
	}

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

	DescribeTable("FromFloat", func(v float64, scale int, exp string) {
		d := decimal.FromFloat(v, uint8(scale))
		Ω(d.String()).Should(Equal(exp))
	},
		Entry("Zero", 0.0, 0, "0"),
		Entry("One", 1.0, 0, "1"),
		Entry("Round", 1.33333, 2, "1.33"),
		Entry("Round up", 1.66666, 2, "1.66"),
		Entry("negative", -1.44444, 2, "-1.44"),
		Entry("scale up", 1.0, 2, "1.00"),
	)

	It("GoStringer", func() {
		d, err := decimal.FromString("3.30")
		Ω(err).Should(Succeed())
		Ω(d.GoString()).Should(Equal("3.30m"))
	})

	It("Valuer", func() {
		d, err := decimal.FromString("3.33")
		Ω(err).Should(Succeed())

		v, err := d.Value()
		Ω(v).Should(Equal("3.33"))
		Ω(err).Should(Succeed())
	})

	DescribeTable("ShortString", func(str, exp string) {
		d, err := decimal.FromString(str)
		Ω(err).Should(Succeed())
		Ω(d.ShortString()).Should(Equal(exp))
	},
		Entry("No ending zero", "3.4", "3.4"),
		Entry("Ending integer zeros", "300", "300"),
		Entry("Ending fragment zeros", "3.00", "3"),
		Entry("Zero", "0.00", "0"),
		Entry("Ten", "10.00", "10"),
		Entry("Ten 2", "10", "10"),
		Entry("3.30", "3.30", "3.3"),
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

		Ω(d.Int64()).Should(Equal(exp))
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
			Ω(d.Neg().String()).Should(Equal(exp))
		},
			Entry("Zero", "0", "0"),
			Entry("Zero 2", "0.00", "0.00"),
			Entry("Positive", "3.456", "-3.456"),
			Entry("Negative", "-3.4456", "3.4456"),
		)

		DescribeTable("Add", func(a, b, c string) {
			assertBinOp(a, b, c, func(x, y decimal.Decimal) interface{} {
				return x.Add(y)
			})
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
			assertBinOp(a, b, c, func(x, y decimal.Decimal) interface{} {
				return x.AddToScale(y, scale)
			})
		},
			Entry("Scale not changed", "3.10", "4.01", "7.11", 2),
			Entry("Extend scale", "3.10", "1", "4.100", 3),
			Entry("Shrink scale round up", "3.1", "0.05", "3.2", 1),
			Entry("Shrink scale round down", "3.1", "0.04", "3.1", 1),
		)

		DescribeTable("Subtract", func(a, b, c string) {
			assertBinOp(a, b, c, func(x, y decimal.Decimal) interface{} {
				return x.Sub(y)
			})
		},
			Entry("Two positive integers", "4", "3", "1"),
			Entry("Positive negative integers", "-4", "3", "-7"),
			Entry("Scale equaled", "2.2", "1.1", "1.1"),
			Entry("Scale not equal", "4.5", "1", "3.5"),
			Entry("Add up to integer", "2.8", "1.2", "1.6"),
			Entry("fragment", "0.0007", "0.0003", "0.0004"),
			Entry("fragment and integer", "300", "0.3", "299.7"),
		)

		DescribeTable("SubtractToScale", func(a, b, c string, scale int) {
			assertBinOp(a, b, c, func(x, y decimal.Decimal) interface{} {
				return x.SubToScale(y, scale)
			})
		},
			Entry("Scale matched", "1.2", "0.2", "1.0", 1),
			Entry("Extend scale", "1.2", "0.2", "1.00", 2),
			Entry("Shrink scale round up", "1", "0.4", "1", 0),
			Entry("Shrink scale round down", "1", "0.6", "0", 0),
		)

		DescribeTable("Multiply", func(a, b, c string) {
			assertBinOp(a, b, c, func(x, y decimal.Decimal) interface{} {
				return x.Mul(y)
			})
		},
			Entry("zero", "3", "0", "0"),
			Entry("Integers", "3", "4", "12"),
			Entry("Integer and fragment", "3", "0.4", "1.2"),
			Entry("fragments", "0.3", "0.5", "0.2"),
			Entry("negative fragement", "-0.3", "0.5", "-0.2"),
			Entry("fragments round down", "0.3", "0.4", "0.1"),
			Entry("big and little 1", "-100", "0.01", "-1.00"),
			Entry("little and little", "-0.001", "-0.01", "0.000"),
			Entry("integer round up", "40", "50", "2000"),
		)

		DescribeTable("MultiplyToScale", func(a, b, c string, scale int) {
			assertBinOp(a, b, c, func(x, y decimal.Decimal) interface{} {
				return x.MulToScale(y, scale)
			})
		},
			Entry("Scale not changed", "1.2", "0.6", "0.7", 1),
			Entry("extend scale", "1.2", "0.6", "0.72", 2),
			Entry("shrink scale round up", "1.5555", "1", "1.6", 1),
			Entry("shrink scale round down", "1.455", "1", "1", 0),
			Entry("shrink scale round down 2", "1.445", "0.1", "0.14", 2),
			Entry("extend scale more", "1.0", "2", "2.000", 3),
		)

		DescribeTable("Div", func(a, b, c string) {
			assertBinOp(a, b, c, func(x, y decimal.Decimal) interface{} {
				return x.Div(y)
			})
		},
			Entry("integer div", "9", "3", "3"),
			Entry("integer div 2", "1", "4", "0"),
			Entry("integer div to fragment round up", "1", "4.0", "0.3"),
			Entry("integer div to fragment round down", "1", "3.0", "0.3"),
			Entry("div with nearest number", "1", "1.00001", "0.99999"),
			Entry("div with small number", "100000", "3", "33333"),
			Entry("div with large number 1", "5", "6000", "0"),
			Entry("div with large number 2", "5.0", "6000", "0.0"),
			Entry("div with large number 3", "5.0000", "6000", "0.0008"),
			Entry("negative div to zero", "-1", "3", "0"),
			Entry("negative div to fragment", "-1", "6.0", "-0.2"),
			Entry("div means multiply", "1", "0.25", "4.00"),
		)

		It("div by zero", func() {
			Ω(func() {
				decimal.FromInt(1).Div(decimal.FromInt(0))
			}).Should(Panic())
		})

		DescribeTable("DivToScale", func(a, b, c string, scale int) {
			assertBinOp(a, b, c, func(x, y decimal.Decimal) interface{} {
				return x.DivToScale(y, scale)
			})
		},
			Entry("scale not changed", "9", "3", "3", 0),
			Entry("expand scale", "10", "4", "2.50", 2),
			Entry("shrink scale round up", "1.454", "1", "1.5", 1),
			Entry("shrink scale round down", "1.454", "1", "1.45", 2),
		)

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

	Context("GetZero", func() {
		for i := 0; i < 9; i++ {
			It(strconv.Itoa(i), func() {
				d := decimal.Zero(i)
				Ω(d.EQ(decimal.FromInt(0))).Should(BeTrue())
				Ω(d.Scale()).Should(Equal(uint8(i)))
			})
		}
	})

	DescribeTable("IsZero", func(s string, is bool) {
		d, err := decimal.FromString(s)
		Ω(err).Should(Succeed())

		Ω(d.IsZero()).Should(Equal(is))
	},
		Entry("is", "0.00", true),
		Entry("not", "0.01", false),
	)

	Context("Compare", func() {
		DescribeTable("Sign", func(s string, sign int) {
			d, err := decimal.FromString(s)
			Ω(err).Should(Succeed())

			Ω(d.Sign()).Should(Equal(sign))
		},
			Entry("Zero", "0", 0),
			Entry("Positive", "1.333", 1),
			Entry("Negative", "-2.3", -1),
		)

		DescribeTable("Cmp", func(a, b string, r int) {
			assertBinOp(a, b, r, func(x, y decimal.Decimal) interface{} {
				return x.Cmp(y)
			})
		},
			Entry("Equal has the same scale", "0", "0", 0),
			Entry("Equal has different scale", "1.00", "1.000", 0),
			Entry("Less than", "1", "9", -1),
			Entry("Greater than", "2.1", "2", 1),
		)

		DescribeTable("LessThan", func(a, b string, r bool) {
			assertBinOp(a, b, r, func(x, y decimal.Decimal) interface{} {
				return x.LT(y)
			})
		},
			Entry("Equal", "0.0", "0.00", false),
			Entry("less", "0.0", "1.00", true),
			Entry("greater", "0.01", "0.00", false),
		)

		DescribeTable("GreaterThan", func(a, b string, r bool) {
			assertBinOp(a, b, r, func(x, y decimal.Decimal) interface{} {
				return x.GT(y)
			})
		},
			Entry("Equal", "0.0", "0.00", false),
			Entry("less", "0.0", "1.00", false),
			Entry("greater", "0.01", "0.00", true),
		)

		DescribeTable("LessTharOrEqual", func(a, b string, r bool) {
			assertBinOp(a, b, r, func(x, y decimal.Decimal) interface{} {
				return x.LTE(y)
			})
		},
			Entry("Equal", "0.0", "0.00", true),
			Entry("less", "0.0", "1.00", true),
			Entry("greater", "0.01", "0.00", false),
		)

		DescribeTable("GreaterThanOrEqual", func(a, b string, r bool) {
			assertBinOp(a, b, r, func(x, y decimal.Decimal) interface{} {
				return x.GTE(y)
			})
		},
			Entry("Equal", "0.0", "0.00", true),
			Entry("less", "0.0", "1.00", false),
			Entry("greater", "0.01", "0.00", true),
		)

		DescribeTable("Equal", func(a, b string, r bool) {
			assertBinOp(a, b, r, func(x, y decimal.Decimal) interface{} {
				return x.EQ(y)
			})
		},
			Entry("Equal", "0.0", "0.00", true),
			Entry("less", "0.0", "1.00", false),
			Entry("greater", "0.01", "0.00", false),
		)

		DescribeTable("NotEqual", func(a, b string, r bool) {
			assertBinOp(a, b, r, func(x, y decimal.Decimal) interface{} {
				return x.NE(y)
			})
		},
			Entry("Equal", "0.0", "0.00", false),
			Entry("less", "0.0", "1.00", true),
			Entry("greater", "0.01", "0.00", true),
		)
	})

	It("bson marshal", func() {
		// bson.Marshal() expected the value is a document, it will fail
		// if marshal Decimal directly, wrap it inside a struct
		var v, back struct {
			V decimal.Decimal
		}

		Ω(bsonRoundTrip(v, &back)).Should(Succeed())
		Ω(back).Should(Equal(v))

		v.V, _ = decimal.FromString("3.00")
		Ω(bsonRoundTrip(v, &back)).Should(Succeed())
		Ω(back).Should(Equal(v))
	})

	It("Json marshal", func() {
		d, err := decimal.FromString("3.30")
		Ω(err).Should(Succeed())
		Ω(json.Marshal(d)).Should(Equal([]byte("3.30")))

		d = decimal.Zero(0)
		Ω(json.Unmarshal([]byte("3.30"), &d)).Should(Succeed())
	})

	Context("decimal128", func() {

		DescribeTable("ToDecimal128", func(d decimal.Decimal, expLow, expHigh uint64) {
			low, high := d.ToDecimal128()
			Ω(low).Should(Equal(expLow))
			Ω(high).Should(Equal(expHigh))
		},
			Entry("Zero", decimal.Zero(0), uint64(0), uint64(0x3040000000000000)),
			Entry("Zero scale 2", decimal.Zero(2), uint64(0), uint64(0x303c000000000000)),
			Entry("One", decimal.FromInt(1), uint64(1), uint64(0x3040000000000000)),
			Entry("Negtative one", decimal.FromInt(-1), uint64(1), uint64(0xb040000000000000)),
		)

		DescribeTable("FromDecimal128", func(low, high uint64, decStr string) {
			exp, err := decimal.FromString(decStr)
			Ω(err).Should(Succeed())
			d := decimal.FromDecimal128(low, high)
			Ω(d).Should(Equal(exp))
		},
			Entry("digits in range", uint64(0x1234567890123456), uint64(0x3040000000000000), "1311768467284833366"),
			// TODO
			XEntry("too many digits only use low", uint64(0x8234567890123456), uint64(0x303e000000000000), "938221899953276219"),
		)

		DescribeTable("FromDecimal128 out of range", func(low, high uint64) {
			Ω(func() {
				decimal.FromDecimal128(low, high)
			}).Should(Panic())
		},
			Entry("too big", uint64(0x8234567890123456), uint64(0x3040000000000000)),
			Entry("scale too large", uint64(0), uint64(0x3042000000000000)),
			Entry("scale too small", uint64(0), uint64(0x302e000000000000)),
			Entry("Two big uses high", uint64(0), uint64(0x3040010000000000)),
		)

		DescribeTable("Round trip", func(s string) {
			d, err := decimal.FromString(s)
			Ω(err).Should(Succeed())
			low, high := d.ToDecimal128()
			back := decimal.FromDecimal128(low, high)
			Ω(back).Should(Equal(d))
		},
			Entry("Zero", "0"),
			Entry("Zero scale 2", "0.00"),
			Entry("One", "1"),
			Entry("Negative one", "-1"),
		)

		It("Random round trip", func() {
			for i := 0; i < 100; i++ {
				scale := rand.Intn(9)
				d, err := decimal.FromStringWithScale(strconv.FormatFloat(rand.Float64()*float64(rand.Int31()), 'f', scale, 64), scale)
				Ω(err).Should(Succeed())
				low, high := d.ToDecimal128()
				back := decimal.FromDecimal128(low, high)
				Ω(back).Should(Equal(d))
			}
		})
	})

})
