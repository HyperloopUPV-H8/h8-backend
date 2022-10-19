package TransportController

import (
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
)

type TransportController struct {
	source      *gopacket.PacketSource
	connections map[string]*connection
}

// Currently hardcoded, this will change once we find a way to autodetect the interface
const interfaceName = "\\Device\\NPF_loopback"
const connectionCheckDelay = time.Second

func New(ips []string, ports []int) (*TransportController, error) {
	connections := make(map[string]*connection)
	for i, ip := range ips {
		conn, err := newConnection(ip, ports[i])
		if err == nil {
			connections[ip] = conn
		}
	}

	source, err := pcap.OpenLive(interfaceName, ^0, true, pcap.BlockForever)
	if err != nil {
		return nil, err
	}

	controller := &TransportController{
		source:      gopacket.NewPacketSource(source, source.LinkType()),
		connections: connections,
	}

	go controller.checkConnections(connectionCheckDelay)

	return controller, nil
}

func (controller *TransportController) checkConnections(delay time.Duration) {
	for {
		for _, conn := range controller.connections {
			conn.checkAlive()
			conn.tryConnect()
		}

		<-time.After(connectionCheckDelay)
	}
}

func (controller *TransportController) Receive() []byte {
	for {
		if packet, err := controller.source.NextPacket(); err == nil && packet != nil && controller.networkFlowFilter(packet) {
			return packet.ApplicationLayer().Payload()
		}
	}
}

func (controller *TransportController) Send(payload []byte, ip string) {
	conn, ok := controller.connections[ip]
	if ok && conn.isAlive {
		_, err := conn.tcp.Write(payload)
		if err != nil {
			conn.disconnect()
		}
	}
}

func (controller *TransportController) networkFlowFilter(packet gopacket.Packet) bool {
	networkLayer := packet.NetworkLayer()
	flow := networkLayer.NetworkFlow()
	_, okSrc := controller.connections[flow.Src().String()]
	_, okDst := controller.connections[flow.Dst().String()]
	return okSrc && okDst
}
