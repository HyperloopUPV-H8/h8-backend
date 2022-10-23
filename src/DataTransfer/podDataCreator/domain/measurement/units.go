package measurement

import (
	"regexp"
	"strconv"

	"github.com/HyperloopUPV-H8/Backend-H8/DataTransfer/utils"
)

type Units struct {
	podUnits     Unit
	displayUnits Unit
}

func NewUnits(podUnitsStr string, displayUnitsStr string) Units {
	podUnits := newUnit(podUnitsStr)
	displayUnits := newUnit(displayUnitsStr)
	return Units{podUnits: podUnits, displayUnits: displayUnits}
}

type Unit struct {
	name       string
	operations []Operation
}

type Operation struct {
	operator string
	operand  float64
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

func getOperations(ops string) []Operation {
	opExp, err := regexp.Compile(`([+\-\/*]{1})(\d+)`)

	if err != nil {
		utils.PrintRegexErr("opExp", err)
	}

	matches := opExp.FindAllStringSubmatch(ops, -1)
	operations := make([]Operation, 0)
	for _, match := range matches {
		operation := getOperation(match[1], match[2])
		operations = append(operations, operation)
	}
	return operations
}

func getOperation(operator string, operand string) Operation {
	num, err := strconv.ParseFloat(operand, 64)

	if err != nil {
		utils.PrintParseNumberErr(err)
	}

	return Operation{operator: operator, operand: num}
}

// func (u Unit) convert(v Value) Value {

// }
