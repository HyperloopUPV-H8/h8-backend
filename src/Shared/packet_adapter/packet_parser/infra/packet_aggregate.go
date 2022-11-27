package infra

import (
	"log"
	"strconv"

	excelAdapter "github.com/HyperloopUPV-H8/Backend-H8/Shared/excel_adapter/domain"
	"github.com/HyperloopUPV-H8/Backend-H8/Shared/packet_adapter/packet_parser/domain"
)

type PacketAggregate struct {
	packets map[uint16]domain.PacketDescriptor
	enums   map[string]domain.Enum
}

func (aggregate PacketAggregate) get(id uint16) domain.PacketDescriptor {
	return aggregate.packets[id]
}

func (aggregate PacketAggregate) getEnum(name string) domain.Enum {
	return aggregate.enums[name]
}

func NewPacketAggregate(boards map[string]excelAdapter.BoardDTO) *PacketAggregate {
	packets := make(map[uint16]domain.PacketDescriptor)
	enums := make(map[string]domain.Enum)

	for _, board := range boards {
		for _, packet := range board.GetPackets() {
			id, err := strconv.ParseInt(packet.Description.ID, 10, 16)
			if err != nil {
				log.Fatalf("packet parser: new packet aggregate: failed to parse id %s: %s\n", packet.Description.ID, err)
			}
			values := make([]domain.ValueDescriptor, len(packet.Measurements))
			for i, value := range packet.Measurements {
				kind := value.ValueType

				if domain.IsEnum(value.ValueType) {
					enums[value.Name] = domain.NewEnum(value.ValueType)
					kind = "enum"
				}

				values[i] = domain.ValueDescriptor{
					Name: value.Name,
					Kind: kind,
				}
			}

			packets[uint16(id)] = values
		}
	}
	return &PacketAggregate{
		packets: packets,
		enums:   enums,
	}
}
