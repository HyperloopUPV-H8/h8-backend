package TransportController

import (
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
)

// TransportController is in charge of sending and receiving messages from/to the pod
type TransportController struct {
	source      *gopacket.PacketSource
	connections map[string]*connection
}

const interfaceName = "\\Device\\NPF_loopback"
const connectionCheckDelay = time.Second

// Create a new TransportService instance trying to connect to the given ips through the given ports
// this will return an error if the packet source couldn't be created. Any invalid ips will be ignored
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

	ts := &TransportController{
		source:      gopacket.NewPacketSource(source, source.LinkType()),
		connections: connections,
	}

	go ts.checkConnections(connectionCheckDelay)

	return ts, nil
}

// Periodically checks all connections to make sure every one is alive
// intended to be run as a goroutine
func (ts *TransportController) checkConnections(delay time.Duration) {
	for {
		for _, conn := range ts.connections {
			conn.checkAlive()
			conn.tryConnect()
		}

		<-time.After(connectionCheckDelay)
	}
}

// Recieve the next packet that isn't meant to be received/sent by/from the backend
func (ts *TransportController) Receive() []byte {
	for {
		if packet, err := ts.source.NextPacket(); err == nil && packet != nil && ts.networkFlowFilter(packet) {
			return packet.ApplicationLayer().Payload()
		}
	}
}

// Send a message to the given ip
func (ts *TransportController) Send(payload []byte, ip string) {
	conn, ok := ts.connections[ip]
	if ok && conn.isAlive {
		_, err := conn.tcp.Write(payload)
		if err != nil {
			conn.disconnect()
		}
	}
}

// Returns if the packet has been sent between connections
func (ts *TransportController) networkFlowFilter(packet gopacket.Packet) bool {
	networkLayer := packet.NetworkLayer()
	flow := networkLayer.NetworkFlow()
	_, okSrc := ts.connections[flow.Src().String()]
	_, okDst := ts.connections[flow.Dst().String()]
	return okSrc && okDst
}
