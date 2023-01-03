package models

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
)

type Operations []Operation

func (operations Operations) Convert(value float64) float64 {
	result := value
	for _, op := range operations {
		result = op.convert(result)
	}
	return result
}

func (operations Operations) Revert(value float64) float64 {
	result := value
	for i := len(operations) - 1; i >= 0; i-- {
		result = operations[i].revert(result)
	}
	return result
}

type Operation struct {
	Operator string
	Operand  float64
}

func (operation Operation) convert(value float64) float64 {
	switch operation.Operator {
	case "+":
		return value + operation.Operand
	case "-":
		return value - operation.Operand
	case "*":
		return value * operation.Operand
	case "/":
		return value / operation.Operand
	}
	return value
}

func (operation Operation) revert(value float64) float64 {
	switch operation.Operator {
	case "+":
		return value - operation.Operand
	case "-":
		return value + operation.Operand
	case "*":
		return value / operation.Operand
	case "/":
		return value * operation.Operand
	}
	return value
}

const decimalRegex = `[-+]?(\d*\.)?\d+(e[-+]?\d+)?`

var operationExp = regexp.MustCompile(fmt.Sprintf(`([+\-\/*]{1})(%s)`, decimalRegex))

func NewOperations(literal string) Operations {
	matches := operationExp.FindAllStringSubmatch(literal, -1)
	operations := make([]Operation, 0)
	for _, match := range matches {
		operation := getOperation(match[1], match[2])
		operations = append(operations, operation)
	}
	return operations
}
func getOperation(operator string, operand string) Operation {
	numOperand, err := strconv.ParseFloat(operand, 64)
	if err != nil {
		log.Fatalln("units: operations: getOperation:", err)
	}
	return Operation{
		Operator: operator,
		Operand:  numOperand,
	}
}
