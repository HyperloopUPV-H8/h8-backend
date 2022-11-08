package value

import (
	"fmt"
	"testing"
)

func TestString(t *testing.T) {
	newEnumValue := "FORWARD"
	specialString := NewDefault("ENUM(FORWARD,BACKWARD)", "cdeg#/100#", "cdeg#/100#")
	specialString.Update(newEnumValue)
	fmt.Println(specialString.ToDisplayUnitsString())
}
