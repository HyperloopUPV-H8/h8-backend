package sniffer

import (
	"encoding/binary"
	"errors"
	"time"

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

type SnifferConfig struct {
	Mtu              uint
	SnifferInterface string `toml:"sniffer_interface"`
}

func New(filter string, config Config) (*Sniffer, error) {
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

func (sniffer *Sniffer) Listen(output chan<- packet.Raw) <-chan error {
	sniffer.trace.Info().Msg("start listening")
	errorChan := make(chan error)

	go sniffer.read(output, errorChan)

	return errorChan
}

func (sniffer *Sniffer) read(output chan<- packet.Raw, errorChan chan<- error) {
	for {
		raw, _, err := sniffer.source.ReadPacketData()
		if err != nil {
			sniffer.trace.Error().Stack().Err(err).Msg("")
			errorChan <- err
			close(errorChan)
			return
		}

		sniffer.trace.Trace().Msg("read")

		packet := gopacket.NewPacket(raw, sniffer.source.LinkType(), gopacket.DecodeOptions{
			NoCopy: true,
		})

		rawPacket, err := sniffer.parseLayers(packet.Layers())
		if err != nil {
			sniffer.trace.Error().Stack().Err(err).Msg("")
			errorChan <- err
			continue
		}

		sniffer.trace.Trace().Msg("parsed")
		output <- rawPacket
	}
}

var syntheticSeqNum uint32 = 0

func (sniffer *Sniffer) parseLayers(packetLayers []gopacket.Layer) (packet.Raw, error) {
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
		return packet.Raw{}, errors.New("failed to get flow")
	}

	//Config endianess from config.toml
	id := binary.LittleEndian.Uint16(payload[:2])

	return packet.Raw{
		Metadata: packet.NewMetaData(from, to, id, seqNum, timestamp),
		Payload:  payload[2:],
	}, nil
}
