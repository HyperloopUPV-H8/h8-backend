package idExpander

import (
	"regexp"
	"strconv"

	"github.com/HyperloopUPV-H8/Backend-H8/DataTransfer/utils"
)

func GetAllIds(id string) []string {
	if isNumber(id) {
		return []string{id}
	}

	prefix, begin, end := getIdParts(id)
	idRange := getIdRange(begin, end)
	finalIds := make([]string, len(idRange))
	for index, sufix := range idRange {
		finalIds[index] = prefix + sufix
	}
	return finalIds
}

func isNumber(id string) bool {
	numberExp, err := regexp.Compile(`^\d+$`)

	if err != nil {
		utils.PrintRegexErr("numberExp", err)
	}

	return numberExp.MatchString(id)
}

func getIdParts(id string) (prefix string, begin string, end string) {
	expandExp, err := regexp.Compile(`^(\d*)\[(\d+),(\d+)\]$`)

	if err != nil {
		utils.PrintRegexErr("expandExp", err)
	}

	matches := expandExp.FindStringSubmatch(id)
	return matches[1], matches[2], matches[3]
}

func getIdRange(begin string, end string) []string {
	lowerLimit := stringToInt(begin)
	upperLimit := stringToInt(end)
	numRange := getRange(lowerLimit, upperLimit)
	stringRange := rangeToString(numRange)
	return stringRange
}

func stringToInt(num string) int {
	n, err := strconv.ParseInt(num, 10, 32)

	if err != nil {
		utils.PrintParseNumberErr(err)
	}

	return int(n)
}

func getRange(n1 int, n2 int) []int {
	numRange := make([]int, n2-n1+1)
	for n := n1; n <= n2; n++ {
		index := n - n1
		numRange[index] = n
	}
	return numRange
}

func rangeToString(numRange []int) []string {
	rangeInString := make([]string, len(numRange))
	for index, num := range numRange {
		str := strconv.Itoa(num)
		rangeInString[index] = str
	}
	return rangeInString
}
