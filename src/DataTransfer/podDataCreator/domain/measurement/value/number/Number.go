package number

import "fmt"

type Number struct {
	value        float32
	podUnits     Unit
	displayUnits Unit
}

func NewNumber(podUnitString string, displayUnitString string) *Number {
	podUnits := newUnit(podUnitString)
	displayUnits := newUnit(displayUnitString)
	return &Number{value: 0, podUnits: podUnits, displayUnits: displayUnits}
}

func (i *Number) ToDisplayString() string {
	number := float32(i.value)
	internationalSystemNumber := undoUnits(number, i.podUnits.operations)
	result := convertToUnits(internationalSystemNumber, i.displayUnits.operations)
	return fmt.Sprintf("%v", result)
}
