package measurement

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

const decimalRegex = `[-+]?(\d*\.)?\d+(e[-+]?\d+)?`

type Ranges struct {
	safe    [2]float64
	warning [2]float64
}

var rangesExp = regexp.MustCompile(fmt.Sprintf(`^\[((%s)*)\,((%s)*)\]$`, decimalRegex, decimalRegex))

func NewRanges(safeRangeStr string, warningRangeStr string) Ranges {
	safeRange := getRangesFromString(strings.ReplaceAll(safeRangeStr, " ", ""))
	warningRange := getRangesFromString(strings.ReplaceAll(warningRangeStr, " ", ""))

	return Ranges{safe: safeRange, warning: warningRange}
}

func getRangesFromString(str string) [2]float64 {
	matches := rangesExp.FindStringSubmatch(str)

	begin := getFloat(matches[1])
	end := getFloat(matches[2])

	return [2]float64{begin, end}
}

func getFloat(intString string) float64 {
	if len(intString) == 0 {
		return 0
	}

	number, err := strconv.ParseFloat(intString, 64)

	if err != nil {
		fmt.Printf("error parsing int: %v\n", err)
	}
	return number
}
