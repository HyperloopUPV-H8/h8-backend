package measurement

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type Ranges struct {
	safe    [2]int
	warning [2]int
}

var rangesExp = regexp.MustCompile(`^\[(-?\d*)\,(-?\d*)\]$`)

func NewRanges(safeRangeStr string, warningRangeStr string) Ranges {
	safeRange := getRangesFromString(strings.ReplaceAll(safeRangeStr, " ", ""))
	warningRange := getRangesFromString(strings.ReplaceAll(warningRangeStr, " ", ""))

	return Ranges{safe: safeRange, warning: warningRange}
}

func getRangesFromString(str string) [2]int {
	matches := rangesExp.FindStringSubmatch(str)

	begin := getInt(matches[1])
	end := getInt(matches[2])

	return [2]int{int(begin), int(end)}
}

func getInt(intString string) int {
	if len(intString) == 0 {
		return 0
	}

	number, err := strconv.ParseInt(intString, 10, 64)

	if err != nil {
		fmt.Printf("error parsing int: %v\n", err)
	}
	return int(number)
}
