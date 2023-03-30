package sniffer

import (
	"os"
	"strconv"

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

func New(dev string, filter string) (*Sniffer, error) {
	trace.Info().Msg("new sniffer")

	source, err := obtainSource(dev, filter)
	if err != nil {
		trace.Error().Stack().Err(err).Msg("")
		return nil, err
	}

	return &Sniffer{
		source: source,
		filter: filter,
		trace:  trace.With().Str("component", "sniffer").Str("dev", dev).Logger(),
	}, nil
}

func obtainSource(dev string, filter string) (*pcap.Handle, error) {
	trace.Debug().Str("dev", dev).Str("filter", filter).Msg("obtain source")
	snaplen, err := strconv.ParseInt(os.Getenv("INTERFACE_MTU"), 10, 32)
	if err != nil {
		trace.Fatal().Stack().Err(err).Str("INTERFACE_MTU", os.Getenv("INTERFACE_MTU")).Msg("")
		return nil, err
	}

	source, err := pcap.OpenLive(dev, int32(snaplen), true, pcap.BlockForever)
	if err != nil {
		return nil, err
	}

	if err := source.SetBPFFilter(filter); err != nil {
		return nil, err
	}

	return source, nil
}

func (sniffer *Sniffer) Listen(output chan<- []byte) <-chan error {
	sniffer.trace.Info().Msg("start listening")
	errorChan := make(chan error)

	go sniffer.read(output, errorChan)

	return errorChan
}

func (sniffer *Sniffer) read(output chan<- []byte, errorChan chan<- error) {
	for {
		raw, _, err := sniffer.source.ReadPacketData()
		if err != nil {
			sniffer.trace.Error().Stack().Err(err).Msg("")
			errorChan <- err
			close(errorChan)
			return
		}

		sniffer.trace.Trace().Msg("read")

		output <- gopacket.NewPacket(raw, layers.LayerTypeEthernet, gopacket.DecodeOptions{
			Lazy:   true,
			NoCopy: true,
		}).ApplicationLayer().Payload()
	}
}
