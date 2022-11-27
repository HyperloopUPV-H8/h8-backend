package packet_adapter

import (
	excelAdapterDomain "github.com/HyperloopUPV-H8/Backend-H8/Shared/excel_adapter/domain"
	loadBalancerInfra "github.com/HyperloopUPV-H8/Backend-H8/Shared/packet_adapter/load_balancer/infra"
	packetParserInfra "github.com/HyperloopUPV-H8/Backend-H8/Shared/packet_adapter/packet_parser/infra"
	transportControllerInfra "github.com/HyperloopUPV-H8/Backend-H8/Shared/packet_adapter/transport_controller/infra"
	transportControllerSniffer "github.com/HyperloopUPV-H8/Backend-H8/Shared/packet_adapter/transport_controller/infra/sniffer"
	unitsInfra "github.com/HyperloopUPV-H8/Backend-H8/Shared/packet_adapter/units/infra"
)

type PacketAdapter struct {
	transportController  *transportControllerInfra.TransportController
	loadBalancer         *loadBalancerInfra.LoadBalancer
	packetAggregate      *packetParserInfra.PacketAggregate
	podUnitAggregate     *unitsInfra.UnitAggregate
	displayUnitAggregate *unitsInfra.UnitAggregate
}

func New(tcSettigns transportControllerInfra.Config, routines int, boards map[string]excelAdapterDomain.BoardDTO) *PacketAdapter {

	return &PacketAdapter{
		transportController:  transportControllerInfra.NewTransportController(tcSettigns),
		packetAggregate:      packetParserInfra.NewPacketAggregate(boards),
		podUnitAggregate:     unitsInfra.NewPodUnitAggregate(boards),
		displayUnitAggregate: unitsInfra.NewDisplayUnitAggregate(boards),
	}
}

type AdapterSettings struct {
	Device        string
	Live          bool
	Port          uint16
	RemoteIPs     []string
	RemotePorts   []uint16
	TCPSnaplen    int32
	SnifferConfig transportControllerSniffer.Config
}
