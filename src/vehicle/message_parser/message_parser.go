package message_parser

import (
	"encoding/json"
	"fmt"

	"github.com/HyperloopUPV-H8/Backend-H8/vehicle/models"
	"github.com/rs/zerolog"
)

type MessageParser struct {
	infoId        uint16
	warningId     uint16
	faultId       uint16
	errorId       uint16
	boardIdToName map[uint]string
	trace         zerolog.Logger
}

func (parser *MessageParser) Parse(id uint16, raw []byte) (any, error) {
	kind, err := parser.getKind(id)

	if err != nil {
		parser.trace.Error().Err(err).Msg("error getting kind")
		return models.ProtectionMessage{}, err
	}

	payload := raw[2:]

	if kind == "info" {
		return parser.toInfoMessage(kind, payload)
	}

	return parser.toProtectionMessage(kind, payload)

}

func (parser *MessageParser) toInfoMessage(kind string, payload []byte) (models.InfoMessage, error) {
	var adapter InfoMessageAdapter
	err := json.Unmarshal(payload, &adapter)

	if err != nil {
		parser.trace.Error().Err(err).Str("message", string(payload)).Msg("error parsing info message")
		return models.InfoMessage{}, err
	}

	name, ok := parser.boardIdToName[adapter.BoardId]

	if !ok {
		parser.trace.Error().Uint("board id", adapter.BoardId).Msg("board id not found")
		name = "DEFAULT"
	}

	return models.InfoMessage{
		Board:     name,
		Timestamp: adapter.Timestamp,
		Msg:       adapter.Msg,
		Kind:      "info",
	}, nil

}

func (parser *MessageParser) toProtectionMessage(kind string, payload []byte) (models.ProtectionMessage, error) {
	var adapter ProtectionMessageAdapter
	err := json.Unmarshal(payload, &adapter)

	if err != nil {
		parser.trace.Error().Err(err).Str("message", string(payload)).Msg("error parsing protection message")
		return models.ProtectionMessage{}, err
	}

	protection, err := getProtection(adapter.Protection.Type, *adapter.Protection.Data)

	if err != nil {
		parser.trace.Error().Err(err).Msg("protection unmarshal failed")
	}

	name, ok := parser.boardIdToName[adapter.BoardId]

	if !ok {
		parser.trace.Error().Uint("board id", adapter.BoardId).Msg("board id not found")
		name = "DEFAULT"
	}

	return models.ProtectionMessage{
		Kind:       kind,
		Board:      name,
		Name:       adapter.Protection.Name,
		Timestamp:  adapter.Timestamp,
		Protection: protection,
	}, nil
}

func (parser *MessageParser) getKind(id uint16) (string, error) {
	if id == parser.faultId {
		return "fault", nil
	}

	if id == parser.warningId {
		return "warning", nil
	}

	if id == parser.errorId {
		return "error", nil
	}

	if id == parser.infoId {
		return "info", nil
	}

	parser.trace.Error().Uint16("id", id).Msg("unrecognized message id")
	return "", fmt.Errorf("unrecognized message id")

}

func getProtection(kind string, payload []byte) (models.Protection, error) {
	switch kind {
	case "OUT_OF_BOUNDS":
		return parseProtection[models.OutOfBounds](kind, payload)
	case "UPPER_BOUND":
		return parseProtection[models.UpperBound](kind, payload)
	case "LOWER_BOUND":
		return parseProtection[models.LowerBound](kind, payload)
	case "EQUALS":
		return parseProtection[models.Equals](kind, payload)
	case "NOT_EQUALS":
		return parseProtection[models.NotEquals](kind, payload)
	case "TIME_ACCUMULATION":
		return parseProtection[models.TimeLimit](kind, payload)
	case "ERROR_HANDLER":
		return parseProtection[models.Error](kind, payload)
	default:
		return models.Protection{}, fmt.Errorf("protection kind not recognized: %s", kind)
	}
}
func parseProtection[T any](kind string, payload []byte) (models.Protection, error) {
	var data T
	err := json.Unmarshal(payload, &data)
	return models.Protection{
		Kind: kind,
		Data: data,
	}, err
}
