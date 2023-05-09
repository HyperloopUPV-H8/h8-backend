package protection_parser

import (
	"encoding/json"
	"fmt"

	"github.com/HyperloopUPV-H8/Backend-H8/vehicle/models"
	"github.com/rs/zerolog"
)

type ProtectionParser struct {
	faultId       uint16
	warningId     uint16
	errorId       uint16
	boardIdToName map[uint]string
	trace         zerolog.Logger
}

func (parser *ProtectionParser) Parse(id uint16, raw []byte) (models.ProtectionMessage, error) {
	kind, err := parser.getKind(id)

	if err != nil {
		parser.trace.Error().Err(err).Msg("error getting kind")
		return models.ProtectionMessage{}, err
	}

	var adapter ProtectionMessageAdapter
	err = json.Unmarshal(raw, &adapter)

	if err != nil {
		parser.trace.Error().Err(err).Str("message", string(raw)).Msg("error parsing protection message")
		return models.ProtectionMessage{}, err
	}

	dataContainer, err := getDataContainer(adapter.Protection.Type)

	if err != nil {
		parser.trace.Error().Err(err).Msg("data container not found")
	}

	err = json.Unmarshal(*adapter.Protection.Data, &dataContainer)

	if err != nil {
		parser.trace.Error().Err(err).Msg("cannot unmarshal protection data")
	}

	name, ok := parser.boardIdToName[adapter.BoardId]

	if !ok {
		parser.trace.Error().Uint("board id", adapter.BoardId).Msg("board id not found")
		return models.ProtectionMessage{}, fmt.Errorf("board id not found: %d", adapter.BoardId)
	}

	return models.ProtectionMessage{
		Kind:      kind,
		Board:     name,
		Name:      adapter.Protection.Name,
		Timestamp: adapter.Timestamp,
		Protection: models.Protection{
			Kind: adapter.Protection.Type,
			Data: dataContainer,
		},
	}, nil
}

func (parser *ProtectionParser) getKind(id uint16) (string, error) {
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

func getDataContainer(kind string) (any, error) {
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
	case "ERROR_HANDLER":
		return models.Error{}, nil
	case "TIME_ACCUMULATION":
		return models.TimeAccumulation{}, nil
	default:
		return nil, fmt.Errorf("protection kind not recognized: %s", kind)
	}
}