package decimal_test

import (
	"encoding/json"

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

	It("Stringer", func() {
		Ω(NullDecimal{}.String()).Should(Equal(""))
		Ω(NullDecimal{FromInt(3), true}.String()).Should(Equal("3"))
	})

	It("Json marshal", func() {
		d, err := FromString("3.30")
		Ω(err).Should(Succeed())
		Ω(json.Marshal(NullDecimal{d, true})).Should(BeEquivalentTo("3.30"))
		Ω(json.Marshal(NullDecimal{d, false})).Should(BeEquivalentTo("null"))

		v := NullDecimal{FromInt(100), false}
		Ω(json.Unmarshal([]byte("3.30"), &v)).Should(Succeed())
		Ω(v).Should(Equal(NullDecimal{d, true}))

		Ω(json.Unmarshal([]byte("null"), &v)).Should(Succeed())
		Ω(v).Should(Equal(NullDecimal{Zero(0), false}))
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
