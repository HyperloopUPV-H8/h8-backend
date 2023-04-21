package parsers

import (
	"bytes"
	"encoding/binary"
	"strconv"
	"strings"

	"github.com/HyperloopUPV-H8/Backend-H8/common"
	excelAdapterModels "github.com/HyperloopUPV-H8/Backend-H8/excel_adapter/models"
	"github.com/HyperloopUPV-H8/Backend-H8/vehicle/models"
)

type ProtectionParser struct {
	Ids       common.Set[uint16]
	faultId   uint16
	warningId uint16
}

type ProtectionParserConfig struct {
	FaultIdKey   string `toml:"fault_id_key"`
	WarningIdKey string `toml:"warning_id_key"`
}

func NewProtectionParser(globalInfo excelAdapterModels.GlobalInfo, config ProtectionParserConfig) ProtectionParser {
	faultId := getId(globalInfo.MessageToId, config.FaultIdKey)
	warningId := getId(globalInfo.MessageToId, config.WarningIdKey)

	ids := common.NewSet[uint16]()

	ids.Add(faultId)
	ids.Add(warningId)

	return ProtectionParser{
		ids,
		faultId,
		warningId,
	}
}

func getId(kindToId map[string]string, key string) uint16 {
	idStr, ok := kindToId[key]

	if !ok {
		//TODO: trace
	}

	id, err := strconv.ParseUint(idStr, 10, 16)

	if err != nil {
		//TODO: trace
	}

	return uint16(id)
}

func (parser *ProtectionParser) Parse(id uint16, raw []byte) (models.Protection, error) {
	parts := strings.Split(string(raw), "\n")

	violation, err := parseViolation(parts[2:])
	if err != nil {
		return models.Protection{}, err
	}

	timestamp, err := parseTimestamp([]byte(parts[len(parts)-1]))
	if err != nil {
		return models.Protection{}, err
	}

	var kind string
	if id == parser.faultId {
		kind = "fault"
	} else {
		kind = "warning"

	}

	return models.Protection{
		Kind:      kind,
		Board:     parts[0],
		Value:     parts[1],
		Violation: violation,
		Timestamp: timestamp,
	}, nil
}

func parseTimestamp(data []byte) (models.Timestamp, error) {
	var timestamp models.Timestamp
	err := binary.Read(bytes.NewReader(data), binary.LittleEndian, &timestamp)
	return timestamp, err
}

var violationStrategy = map[string]func([]string) (models.Violation, error){
	"OUT_OF_BOUNDS": parseOutOfBounds,
	"UPPER_BOUND":   parseUpperBound,
	"LOWER_BOUND":   parseLowerBound,
	"EQUALS":        parseEquals,
	"NOT_EQUALS":    parseNotEquals,
}

func parseViolation(data []string) (models.Violation, error) {
	return violationStrategy[data[0]](data[1:])
}

func parseOutOfBounds(parts []string) (models.Violation, error) {
	violation := models.OutOfBoundsViolation{
		Kind: "OUT_OF_BOUNDS",
	}
	var err error

	violation.Got, err = strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return violation, err
	}

	violation.Want[0], err = strconv.ParseFloat(parts[1], 64)
	if err != nil {
		return violation, err
	}

	violation.Want[1], err = strconv.ParseFloat(parts[2], 64)
	if err != nil {
		return violation, err
	}

	return violation, nil
}

func parseUpperBound(parts []string) (models.Violation, error) {
	violation := models.UpperBoundViolation{
		Kind: "UPPER_BOUND",
	}
	var err error

	violation.Got, err = strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return violation, err
	}

	violation.Want, err = strconv.ParseFloat(parts[1], 64)
	if err != nil {
		return violation, err
	}

	return violation, nil
}

func parseLowerBound(parts []string) (models.Violation, error) {
	violation := models.LowerBoundViolation{
		Kind: "LOWER_BOUND",
	}
	var err error

	violation.Got, err = strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return violation, err
	}

	violation.Want, err = strconv.ParseFloat(parts[1], 64)
	if err != nil {
		return violation, err
	}

	return violation, nil
}

func parseEquals(parts []string) (models.Violation, error) {
	violation := models.EqualsViolation{
		Kind: "EQUALS",
	}
	var err error

	violation.Got, err = strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return violation, err
	}

	return violation, nil
}

func parseNotEquals(parts []string) (models.Violation, error) {
	violation := models.NotEqualsViolation{
		Kind: "NOT_EQUALS",
	}
	var err error

	violation.Got, err = strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return violation, err
	}

	violation.Want, err = strconv.ParseFloat(parts[1], 64)
	if err != nil {
		return violation, err
	}

	return violation, nil
}
