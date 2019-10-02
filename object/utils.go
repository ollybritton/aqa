package object

import (
	"fmt"
	"math"
	"strconv"
)

// IntegerToFloat converts an integer object to a float object.
func IntegerToFloat(i *Integer) *Float {
	return &Float{Value: float64(i.Value)}
}

// IntegerToString converts an integer to a string.
func IntegerToString(i *Integer) *String {
	val := fmt.Sprintf("%d", i.Value)
	return &String{Value: val}
}

// FloatToInteger converts a flaot object to an integer object.
func FloatToInteger(f *Float) *Integer {
	val := int64(math.Round(f.Value))
	return &Integer{Value: val}
}

// FloatToString converts a float into a string.
func FloatToString(f *Float) *String {
	val := fmt.Sprintf("%f", f.Value)
	return &String{Value: val}
}

// StringToInteger converts a string to an integer.
func StringToInteger(s *String) (*Integer, error) {
	val, err := strconv.ParseInt(s.Value, 0, 64)
	if err != nil {
		return &Integer{}, err
	}

	return &Integer{Value: val}, nil
}

// StringToFloat converts a string to an integer.
func StringToFloat(s *String) (*Float, error) {
	val, err := strconv.ParseFloat(s.Value, 10)
	if err != nil {
		return &Float{}, err
	}

	return &Float{Value: val}, nil
}
