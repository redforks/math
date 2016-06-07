package decimal

import "database/sql/driver"

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
