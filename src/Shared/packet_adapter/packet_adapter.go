package packet_adapter

import (
	"os"

	excelAdapterDomain "github.com/HyperloopUPV-H8/Backend-H8/Shared/excel_adapter/domain"
	loadBalancerInfra "github.com/HyperloopUPV-H8/Backend-H8/Shared/packet_adapter/load_balancer/infra"
	packetParserInfra "github.com/HyperloopUPV-H8/Backend-H8/Shared/packet_adapter/packet_parser/infra"
	packetParserDTO "github.com/HyperloopUPV-H8/Backend-H8/Shared/packet_adapter/packet_parser/infra/dto"
	transportControllerInfra "github.com/HyperloopUPV-H8/Backend-H8/Shared/packet_adapter/transport_controller/infra"
	unitsInfra "github.com/HyperloopUPV-H8/Backend-H8/Shared/packet_adapter/units/infra"
	unitsMappers "github.com/HyperloopUPV-H8/Backend-H8/Shared/packet_adapter/units/infra/mappers"
)

type PacketAdapter struct {
	transportController  *transportControllerInfra.TransportController
	loadBalancer         *loadBalancerInfra.LoadBalancer
	packetAggregate      *packetParserInfra.PacketAggregate
	podUnitAggregate     *unitsInfra.UnitAggregate
	displayUnitAggregate *unitsInfra.UnitAggregate
	data                 chan<- packetParserDTO.PacketUpdate
	message              chan<- packetParserDTO.PacketUpdate
	order                <-chan packetParserDTO.PacketUpdate
}

func New(tcSettigns transportControllerInfra.Config, routines int, routineBuf int, data chan<- packetParserDTO.PacketUpdate, message chan<- packetParserDTO.PacketUpdate, order <-chan packetParserDTO.PacketUpdate, boards map[string]excelAdapterDomain.BoardDTO) *PacketAdapter {
	channels := make([]chan []byte, routines)
	for i := 0; i < routines; i++ {
		channels[i] = make(chan []byte, routineBuf)
	}

	txChannels := make([]chan<- []byte, routines)
	for i, channel := range channels {
		txChannels[i] = channel
	}

	adapter := &PacketAdapter{
		transportController:  transportControllerInfra.NewTransportController(tcSettigns),
		packetAggregate:      packetParserInfra.NewPacketAggregate(boards),
		podUnitAggregate:     unitsInfra.NewPodUnitAggregate(boards),
		displayUnitAggregate: unitsInfra.NewDisplayUnitAggregate(boards),
		loadBalancer:         loadBalancerInfra.Init(txChannels),
		data:                 data,
		message:              message,
		order:                order,
	}

	go func(adapter *PacketAdapter) {
		for {
			payload, _ := adapter.transportController.ReceiveData()
			adapter.loadBalancer.Next(payload)
		}
	}(adapter)

	go func(adapter *PacketAdapter) {
		for {
			payload := <-order
			unitsMappers.RevertUpdate(&payload, *adapter.displayUnitAggregate)
			unitsMappers.RevertUpdate(&payload, *adapter.podUnitAggregate)
			adapter.transportController.Send(os.Getenv("TARGET_IP"), packetParserInfra.Encode(payload, *adapter.packetAggregate))
		}
	}(adapter)

	adapter.transportController.OnRead(func(payload []byte) {
		update := packetParserInfra.Decode(payload, *adapter.packetAggregate)
		unitsMappers.ConvertUpdate(&update, *adapter.podUnitAggregate)
		unitsMappers.ConvertUpdate(&update, *adapter.displayUnitAggregate)
		adapter.message <- update
	})

	for _, channel := range channels {
		go func(source <-chan []byte, adapter *PacketAdapter) {
			for {
				payload := <-source
				update := packetParserInfra.Decode(payload, *adapter.packetAggregate)
				unitsMappers.ConvertUpdate(&update, *adapter.podUnitAggregate)
				unitsMappers.ConvertUpdate(&update, *adapter.displayUnitAggregate)
				adapter.data <- update
			}
		}(channel, adapter)
	}

	return adapter
}
