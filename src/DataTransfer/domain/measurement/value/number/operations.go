package number

import (
	"log"
	"regexp"
	"strconv"
)

type Operation struct {
	operator string
	operand  float64
}

var operationExp = regexp.MustCompile(`([+\-\/*]{1})([-+]?(\d*\.)?\d+(e[-+]?\d+)?)`)

func doOperation(number float64, operation Operation) float64 {
	switch operation.operator {
	case "+":
		return number + operation.operand
	case "-":
		return number - operation.operand
	case "*":
		return number * operation.operand
	case "/":
		return number / operation.operand
	default:
		panic("Invalid operation")
	}
}

func getOperations(ops string) []Operation {
	matches := operationExp.FindAllStringSubmatch(ops, -1)
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
		log.Fatalf("parse: %s\n", err)
	}

	return Operation{operator: operator, operand: num}
}

func getOpositeAndReversedOperations(operations []Operation) []Operation {
	newOperations := make([]Operation, len(operations))
	for index, operation := range operations {
		newOperations[index] = getOpositeOperation(operation)
	}
	newOperations = getReversedOperations(newOperations)
	return newOperations
}

func getOpositeOperation(operation Operation) Operation {
	opositeOperation := Operation{operand: operation.operand}
	switch operation.operator {
	case "+":
		opositeOperation.operator = "-"
	case "-":
		opositeOperation.operator = "+"

	case "*":
		opositeOperation.operator = "/"

	case "/":
		opositeOperation.operator = "*"
	default:
		panic("Invalid operator")
	}

	return opositeOperation
}

func getReversedOperations(operations []Operation) []Operation {
	length := len(operations)
	reversedOperations := make([]Operation, length)

	for index, operation := range operations {
		reversedOperations[length-1-index] = operation
	}

	return reversedOperations
}
