package sniffer

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/HyperloopUPV-H8/Backend-H8/common"
	"github.com/HyperloopUPV-H8/Backend-H8/info"
	"github.com/HyperloopUPV-H8/Backend-H8/packet"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/rs/zerolog"
	trace "github.com/rs/zerolog/log"
)

const SNAPLEN = 1500

type Sniffer struct {
	source *pcap.Handle
	filter string
	config Config
	trace  zerolog.Logger
}

func CreateSniffer(info info.Info, config Config, trace zerolog.Logger) Sniffer {
	ips := common.Values(info.Addresses.Boards)
	filter := getFilter(ips, info.Ports.UDP, info.Ports.TcpClient, info.Ports.TcpServer)
	sniffer, err := newSniffer(filter, config)

	if err != nil {
		trace.Fatal().Stack().Err(err).Msg("error creating sniffer")
	}

	return *sniffer
}

func newSniffer(filter string, config Config) (*Sniffer, error) {
	trace.Info().Msg("new sniffer")
	source, err := newSource(config, filter)

	if err != nil {
		trace.Error().Stack().Err(err).Msg("")
		return nil, err
	}

	return &Sniffer{
		source: source,
		filter: filter,
		config: config,
		trace:  trace.With().Str("component", "sniffer").Str("dev", config.Interface).Logger(),
	}, nil
}

func newSource(config Config, filter string) (*pcap.Handle, error) {
	source, err := obtainSource(config.Interface, filter, config.Mtu)

	if err != nil {
		return nil, err
	}

	return source, nil
}

func getFilter(addrs []net.IP, udpPort uint16, tcpClientPort uint16, tcpServerPort uint16) string {
	ipipFilter := getIPIPfilter()
	udpFilter := getUDPFilter(addrs, udpPort)
	tcpFilter := getTCPFilter(addrs, tcpServerPort, tcpClientPort)
	filter := fmt.Sprintf("(%s) or (%s) or (%s)", ipipFilter, udpFilter, tcpFilter)

	trace.Trace().Any("addrs", addrs).Str("filter", filter).Msg("new filter")
	return filter
}

func getIPIPfilter() string {
	return "ip[9] == 4"
}

func getUDPFilter(addrs []net.IP, port uint16) string {
	udp := fmt.Sprintf("(udp port %d)", port)
	udpAddr := ""
	for _, addr := range addrs {
		udpAddr = fmt.Sprintf("%s or (src host %s)", udpAddr, addr.String())
	}
	udpAddr = strings.TrimPrefix(udpAddr, " or ")
	return fmt.Sprintf("%s and (%s)", udp, udpAddr)
}

func getTCPFilter(addrs []net.IP, serverPort uint16, clientPort uint16) string {
	ports := fmt.Sprintf("tcp port %d or %d", serverPort, clientPort)
	flags := "tcp[tcpflags] & (tcp-fin | tcp-syn | tcp-rst) == 0"
	nonZeroPayload := "tcp[tcpflags] & tcp-push != 0"

	srcAddresses := make([]string, 0, len(addrs))
	dstAddresses := make([]string, 0, len(addrs))

	for _, addr := range addrs {
		srcAddresses = append(srcAddresses, fmt.Sprintf("(src host %s)", addr.String()))
		dstAddresses = append(dstAddresses, fmt.Sprintf("(dst host %s)", addr.String()))
	}

	srcAddrsStr := strings.Join(srcAddresses, " or ")
	dstAddrsStr := strings.Join(dstAddresses, " or ")

	filter := fmt.Sprintf("(%s) and (%s) and (%s) and (%s) and (%s)", ports, flags, nonZeroPayload, srcAddrsStr, dstAddrsStr)
	return filter
}

func obtainSource(dev string, filter string, mtu uint) (*pcap.Handle, error) {
	trace.Debug().Str("dev", dev).Str("filter", filter).Msg("obtain source")

	source, err := pcap.OpenLive(dev, int32(mtu), true, pcap.BlockForever)
	if err != nil {
		return nil, err
	}

	if err := source.SetBPFFilter(filter); err != nil {
		return nil, err
	}

	return source, nil
}

func (sniffer *Sniffer) Listen(output chan<- packet.Packet) {
	go sniffer.startReadLoop(output)
}

func (sniffer *Sniffer) startReadLoop(output chan<- packet.Packet) {
	for {
		source, err := newSource(sniffer.config, sniffer.filter)
		if err != nil {
			continue
		}
		sniffer.source = source

		sniffer.trace.Info().Msg("start listening")
		sniffer.read(output)
	}

}

func (sniffer *Sniffer) read(output chan<- packet.Packet) {
	for {
		raw, _, err := sniffer.source.ReadPacketData()
		if err != nil {
			sniffer.trace.Error().Stack().Err(err).Msg("")
			return
		}

		sniffer.trace.Trace().Msg("read")

		packet := gopacket.NewPacket(raw, sniffer.source.LinkType(), gopacket.DecodeOptions{
			NoCopy: true,
		})

		rawPacket, err := sniffer.parseLayers(packet.Layers())
		if err != nil {
			sniffer.trace.Error().Stack().Err(err).Msg("")
			continue
		}

		sniffer.trace.Trace().Msg("parsed")
		output <- rawPacket
	}
}

var syntheticSeqNum uint32 = 0

func (sniffer *Sniffer) parseLayers(packetLayers []gopacket.Layer) (packet.Packet, error) {
	timestamp := time.Now()
	from := ""
	to := ""
	seqNum := syntheticSeqNum
	var payload []byte

layerLoop:
	for _, layer := range packetLayers {
		switch layer := layer.(type) {
		case *layers.IPv4:
			if layer.Protocol == 4 {
				continue layerLoop
			}
			from = layer.SrcIP.String()
			to = layer.DstIP.String()
		case *layers.TCP:
			seqNum = layer.Seq
			payload = layer.Payload
			break layerLoop
		case *layers.UDP:
			syntheticSeqNum++
			payload = layer.Payload
			break layerLoop
		}
	}

	if from == "" || to == "" {
		return packet.Packet{}, errors.New("failed to get flow")
	}

	//Config endianess from config.toml
	if len(payload) < 2 {
		return packet.Packet{}, errors.New("payload smaller than 2")
	}

	id := binary.LittleEndian.Uint16(payload[:2])

	return packet.Packet{
		Metadata: packet.NewMetaData(from, to, id, seqNum, timestamp),
		Payload:  payload[2:],
	}, nil
}
