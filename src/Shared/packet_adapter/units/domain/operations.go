package domain

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
)

type Operations []operation

type operation struct {
	operator string
	operand  float64
}

const decimalRegex = `[-+]?(\d*\.)?\d+(e[-+]?\d+)?`

var operationExp = regexp.MustCompile(fmt.Sprintf(`([+\-\/*]{1})(%s)`, decimalRegex))

func NewOperations(literal string) Operations {
	matches := operationExp.FindAllStringSubmatch(literal, -1)
	operations := make([]operation, 0)
	for _, match := range matches {
		operation := getOperation(match[1], match[2])
		operations = append(operations, operation)
	}
	return operations
}

func getOperation(operator string, operand string) operation {
	numOperand, err := strconv.ParseFloat(operand, 64)
	if err != nil {
		log.Fatalln("units: operations: getOperation:", err)
	}

	return operation{
		operator: operator,
		operand:  numOperand,
	}
}

func DoOperations(operations Operations, value float64) (converted float64) {
	converted = value
	for _, operation := range operations {
		converted = doOperation(operation, converted)
	}
	return converted
}

func doOperation(operation operation, value float64) (converted float64) {
	switch operation.operator {
	case "+":
		return value + operation.operand
	case "-":
		return value - operation.operand
	case "*":
		return value * operation.operand
	case "/":
		return value / operation.operand
	default:
		return value
	}
}

func DoReverseOperations(operations Operations, value float64) (converted float64) {
	converted = value
	for i := (len(operations) - 1); i >= 0; i-- {
		converted = doReverseOperation(operations[i], converted)
	}
	return converted
}

func doReverseOperation(operation operation, value float64) (converted float64) {
	switch operation.operator {
	case "+":
		return value - operation.operand
	case "-":
		return value + operation.operand
	case "*":
		return value / operation.operand
	case "/":
		return value * operation.operand
	default:
		return value
	}
}
