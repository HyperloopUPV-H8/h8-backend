package sniffer

import (
	"encoding/binary"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/HyperloopUPV-H8/Backend-H8/common"
	excel_models "github.com/HyperloopUPV-H8/Backend-H8/excel_adapter/models"
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
	trace  zerolog.Logger
}

func CreateSniffer(global excel_models.GlobalInfo, config Config, trace zerolog.Logger) Sniffer {
	ips := common.Values(global.BoardToIP)
	ips = append(ips, global.BackendIP)
	filter := getFilter(ips, global.ProtocolToPort, config.TcpClientTag, config.TcpServerTag, config.UdpTag)
	sniffer, err := newSniffer(filter, config)

	if err != nil {
		trace.Fatal().Stack().Err(err).Msg("error creating sniffer")
	}

	return *sniffer
}

func newSniffer(filter string, config Config) (*Sniffer, error) {
	trace.Info().Msg("new sniffer")

	source, err := obtainSource(config.Interface, filter, config.Mtu)
	if err != nil {
		trace.Error().Stack().Err(err).Msg("")
		return nil, err
	}

	return &Sniffer{
		source: source,
		filter: filter,
		trace:  trace.With().Str("component", "sniffer").Str("dev", config.Interface).Logger(),
	}, nil
}

func getFilter(addrs []string, protocolToPort map[string]string, tcpClientTag string, tcpServerTag string, udpTag string) string {

	ipipFilter := getIPIPfilter()
	udpFilter := getUDPFilter(addrs, protocolToPort, udpTag)
	tcpFilter := getTCPFilter(addrs, protocolToPort, tcpClientTag, tcpServerTag)
	// tcpFilter := "(tcp port 50500 or 50501) and (src host 127.0.0.9) and (tcp[tcpflags] & (tcp-fin | tcp-syn) == 0)"
	filter := fmt.Sprintf("(%s) or (%s) or (%s)", ipipFilter, udpFilter, tcpFilter)

	trace.Trace().Strs("addrs", addrs).Str("filter", filter).Msg("new filter")
	return filter
}

func getIPIPfilter() string {
	return "ip[9] == 4"
}

func getUDPFilter(addrs []string, protocolToPort map[string]string, udpTag string) string {
	udp := fmt.Sprintf("(udp port %s)", protocolToPort[udpTag])
	udpAddr := ""
	for _, addr := range addrs {
		udpAddr = fmt.Sprintf("%s or (src host %s)", udpAddr, addr)
	}
	udpAddr = strings.TrimPrefix(udpAddr, " or ")
	return fmt.Sprintf("%s and (%s)", udp, udpAddr)
}

func getTCPFilter(addrs []string, protocolToPort map[string]string, tcpClientTag string, tcpServerTag string) string {

	ports := fmt.Sprintf("tcp port %s or %s", protocolToPort[tcpClientTag], protocolToPort[tcpServerTag])
	flags := "tcp[tcpflags] & (tcp-fin | tcp-syn | tcp-rst) == 0"
	nonZeroPayload := "tcp[tcpflags] & tcp-push != 0"

	srcAddresses := make([]string, 0, len(addrs))
	dstAddresses := make([]string, 0, len(addrs))

	for _, addr := range addrs {
		srcAddresses = append(srcAddresses, fmt.Sprintf("(src host %s)", addr))
		dstAddresses = append(dstAddresses, fmt.Sprintf("(dst host %s)", addr))
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
	sniffer.trace.Info().Msg("start listening")

	go sniffer.read(output)

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
