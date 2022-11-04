package number

import "fmt"

type Number struct {
	value        float64
	podUnits     Unit
	displayUnits Unit
}

func NewNumber(podUnitString string, displayUnitString string) *Number {
	podUnits := newUnit(podUnitString)
	displayUnits := newUnit(displayUnitString)
	return &Number{value: 0, podUnits: podUnits, displayUnits: displayUnits}
}

func (n Number) GetPodUnits() string {
	return n.podUnits.name
}

func (n Number) GetDisplayUnits() string {
	return n.displayUnits.name
}

func (n Number) ToPodUnitsString() string {
	number := float64(n.value)
	return fmt.Sprintf("%v", number)
}

func (n Number) ToDisplayUnitsString() string {
	number := float64(n.value)
	internationalSystemNumber := undoUnits(number, n.podUnits.operations)
	result := convertToUnits(internationalSystemNumber, n.displayUnits.operations)
	return fmt.Sprintf("%v", result)
}

func (n *Number) Update(newValue any) {
	newNumber, ok := newValue.(float64)

	if !ok {
		panic("invalid value")
	}
	n.value = newNumber
}
