package number

import (
	"regexp"

	"github.com/HyperloopUPV-H8/Backend-H8/DataTransfer/utils"
)

type Unit struct {
	name       string
	operations []Operation
}

func newUnit(unitStr string) Unit {
	unitExp, err := regexp.Compile(`^([a-zA-Z]+)#((?:[+\-\/*]{1}\d+)+)#$`)

	if err != nil {
		utils.PrintRegexErr("unitExp", err)
	}

	matches := unitExp.FindStringSubmatch(unitStr)
	unit := Unit{
		name:       matches[1],
		operations: getOperations(matches[2]),
	}
	return unit
}

func convertToUnits(number float32, operations []Operation) float32 {
	result := number
	for _, operation := range operations {
		result = doOperation(result, operation)
	}
	return result
}

func undoUnits(number float32, operations []Operation) float32 {
	newOperations := getOpositeAndReversedOperations(operations)
	return convertToUnits(number, newOperations)
}
