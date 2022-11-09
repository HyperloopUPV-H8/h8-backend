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

func (op operation) convert(value float64) float64 {
	switch op.operator {
	case "+":
		return value + op.operand
	case "-":
		return value - op.operand
	case "*":
		return value * op.operand
	case "/":
		return value / op.operand
	}
	return value
}

func (op operation) revert(value float64) float64 {
	switch op.operator {
	case "+":
		return value - op.operand
	case "-":
		return value + op.operand
	case "*":
		return value / op.operand
	case "/":
		return value * op.operand
	}
	return value
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

func (operations Operations) Convert(value any) any {
	final := value.(float64)
	for _, operation := range operations {
		final = operation.convert(final)
	}
	return final
}

func (operations Operations) Revert(value any) any {
	final := value.(float64)
	for _, operation := range operations {
		final = operation.revert(final)
	}
	return final
}
