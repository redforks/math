package decimal

import (
	"database/sql/driver"

	"gopkg.in/mgo.v2/bson"
)

// NullDecimal nullable decimal value
type NullDecimal struct {
	Decimal Decimal
	Valid   bool // Valid is true if current value is not NULL
}

// Value implement driver.Valuer interface.
func (d NullDecimal) Value() (driver.Value, error) {
	if d.Valid {
		return d.Decimal.Value()
	}

	return nil, nil
}

func (d NullDecimal) String() string {
	if d.Valid {
		return d.Decimal.String()
	}
	return ""
}

// GetBSON implement bson.Getter interface, marshal value to mongoDB.
// Marshal to string to pressure both scale and value.
func (d NullDecimal) GetBSON() (interface{}, error) {
	if d.Valid {
		return d.Decimal.GetBSON()
	}

	return bson.Raw{
		Kind: 10,
	}, nil
}

// SetBSON implement bson.Setter interface, marshal value from mongoDB
func (d *NullDecimal) SetBSON(raw bson.Raw) error {
	if raw.Kind == 10 { // 10 means null, see: https://docs.mongodb.com/manual/reference/bson-types/
		d.Valid = false
		return nil
	}

	if err := (&d.Decimal).SetBSON(raw); err != nil {
		return err
	}

	d.Valid = true
	return nil
}

var (
	_ bson.Getter = NullDecimal{}
	_ bson.Setter = &NullDecimal{}
)
