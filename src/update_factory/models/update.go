package models

import "fmt"

type Update struct {
	Id        uint16         `json:"id"`
	HexValue  string         `json:"hexValue"`
	Count     uint64         `json:"count"`
	CycleTime uint64         `json:"cycleTime"`
	Values    map[string]any `json:"measurementUpdates"`
}

type UpdateValue interface {
	Kind() string
}

type NumericValue struct {
	Value   float64 `json:"last"`
	Average float64 `json:"average"`
}

func (numeric NumericValue) Kind() string {
	return "numeric"
}

type BooleanValue bool

func (boolean BooleanValue) Kind() string {
	return "boolean"
}

type EnumValue string

func (enum EnumValue) Kind() string {
	return "enum"
}

type ArrayValue struct {
	Arr any
}

func (arr ArrayValue) Kind() string {
	return "array"
}

func (arr ArrayValue) MarshalJSON() ([]byte, error) {
	if arr.Arr == nil {
		return []byte("null"), nil
	}

	switch typedArr := arr.Arr.(type) {
	case []uint8:
		return getJsonArray[uint8](typedArr), nil
	case []uint16:
		return getJsonArray(typedArr), nil
	case []uint32:
		return getJsonArray(typedArr), nil
	case []uint64:
		return getJsonArray(typedArr), nil
	case []int8:
		return getJsonArray(typedArr), nil
	case []int16:
		return getJsonArray(typedArr), nil
	case []int32:
		return getJsonArray(typedArr), nil
	case []int64:
		return getJsonArray(typedArr), nil
	case []float32:
		return getJsonArray(typedArr), nil
	case []float64:
		return getJsonArray(typedArr), nil
	case []bool:
		return getJsonArray(typedArr), nil
	default:
		return []byte("[]"), nil
	}
}

func getJsonArray[T any](arr []T) []byte {
	arrStr := "["
	for index, item := range arr {
		arrStr = fmt.Sprintf("%s%v", arrStr, item)
		if index < len(arr)-1 {
			arrStr = fmt.Sprintf("%s,", arrStr)
		}
	}
	arrStr = fmt.Sprintf("%s]", arrStr)
	return []byte(arrStr)
}
