package decimal_test

import (
	. "github.com/redforks/math/decimal"
	"gopkg.in/mgo.v2/bson"

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

	Context("bson marshal", func() {
		// bson.Marshal() expected the value is a document, it will fail
		// if marshal NullDecimal directly, wrap it inside a struct
		var v, back struct {
			V NullDecimal
		}

		BeforeEach(func() {
			v.V = NullDecimal{}
			back.V.Valid = true
		})

		It("null", func() {
			Ω(bsonRoundTrip(v, &back)).Should(Succeed())
			Ω(back).Should(Equal(v))
		})

		It("not null", func() {
			v.V = NullDecimal{Zero(1), true}
			Ω(bsonRoundTrip(v, &back)).Should(Succeed())
			Ω(back).Should(Equal(v))
		})
	})

})

// marshal value to bson, then marshal back.
func bsonRoundTrip(v interface{}, back interface{}) error {
	buf, err := bson.Marshal(v)
	if err != nil {
		return err
	}

	return bson.Unmarshal(buf, back)
}
