package util

import "github.com/shopspring/decimal"

// StringPtr returns a pointer to the passed string.
func StringPtr(s string) *string {
	return &s
}

// IntPtr returns a pointer to the passed int.
func IntPtr(s int64) *int64 {
	return &s
}

// FloatPtr returns a pointer to the passed float.
func FloatPtr(s float64) *float64 {
	return &s
}

func FloatToDecimal(s *float64) *decimal.Decimal {
	if s == nil {
		return nil
	} else {
		return DecimalPtr(decimal.NewFromFloat(*s))
	}
}

func DecimalToFloat(s *decimal.Decimal) *float64 {
	if s == nil {
		return nil
	} else {
		return FloatPtr(s.InexactFloat64())
	}
}

// DecimalPtr returns a pointer to the passed decimal.
func DecimalPtr(s decimal.Decimal) *decimal.Decimal {
	return &s
}
