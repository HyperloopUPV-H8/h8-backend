package message_parser

import (
	"encoding/json"
	"fmt"

	"github.com/HyperloopUPV-H8/Backend-H8/vehicle/models"
	"github.com/rs/zerolog"
)

type MessageParser struct {
	faultId       uint16
	warningId     uint16
	errorId       uint16
	boardIdToName map[uint]string
	trace         zerolog.Logger
}

func (parser *MessageParser) Parse(id uint16, raw []byte) (models.ProtectionMessage, error) {
	kind, err := parser.getKind(id)

	if err != nil {
		parser.trace.Error().Err(err).Msg("error getting kind")
		return models.ProtectionMessage{}, err
	}

	payload := raw[2:]

	var adapter MessageAdapter
	err = json.Unmarshal(payload, &adapter)

	if err != nil {
		parser.trace.Error().Err(err).Str("message", string(raw)).Msg("error parsing protection message")
		return models.ProtectionMessage{}, err
	}

	return parser.toProtectionMessage(kind, adapter), nil
}

func (parser *MessageParser) toProtectionMessage(kind string, adapter MessageAdapter) models.ProtectionMessage {
	protectionContainer, err := getProtectionContainer(adapter.Protection.Type)

	if err != nil {
		parser.trace.Error().Err(err).Msg("data container not found")
	}

	err = json.Unmarshal(*adapter.Protection.Data, &protectionContainer)

	if err != nil {
		parser.trace.Error().Err(err).Msg("cannot unmarshal protection data")
	}

	name, ok := parser.boardIdToName[adapter.BoardId]

	if !ok {
		parser.trace.Error().Uint("board id", adapter.BoardId).Msg("board id not found")
		name = "DEFAULT"
	}

	return models.ProtectionMessage{
		Kind:      kind,
		Board:     name,
		Name:      adapter.Protection.Name,
		Timestamp: adapter.Timestamp,
		Protection: models.Protection{
			Kind: adapter.Protection.Type,
			Data: protectionContainer,
		},
	}
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

	parser.trace.Error().Uint16("id", id).Msg("unrecognized message id")
	return "", fmt.Errorf("unrecognized message id")

}

func getProtectionContainer(kind string) (any, error) {
	switch kind {
	case "OUT_OF_BOUNDS":
		return models.OutOfBounds{}, nil
	case "UPPER_BOUND":
		return models.UpperBound{}, nil
	case "LOWER_BOUND":
		return models.LowerBound{}, nil
	case "EQUALS":
		return models.Equals{}, nil
	case "NOT_EQUALS":
		return models.NotEquals{}, nil
	case "TIME_ACCUMULATION":
		return models.TimeAccumulation{}, nil
	case "ERROR_HANDLER":
		return "", nil
	default:
		return nil, fmt.Errorf("protection kind not recognized: %s", kind)
	}
}
