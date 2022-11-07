package streaming

import (
	"time"

	"github.com/HyperloopUPV-H8/Backend-H8/DataTransfer/domain"
	packetParser "github.com/HyperloopUPV-H8/Backend-H8/Shared/packet_adapter/packet_parser/domain"
	"github.com/gorilla/websocket"
)

func DataSocketHandler(ws websocket.Conn, packetChannel chan domain.Packet) {
	go func() {
		for {
			packetWebAdapterBuf := make(map[uint16]PacketWebAdapter, 100)
			timeout := time.After(time.Millisecond * 20)
		loop:
			for {
				select {
				case packet := <-packetChannel:
					adapter := newPacketWebAdapter(packet)
					packetWebAdapterBuf[adapter.Id] = adapter
				case <-timeout:
					ws.WriteJSON(packetWebAdapterBuf)
					break loop
				}
			}
		}
	}()
}

func OrderSocketHandler(ws websocket.Conn, orderWAChannel chan OrderWebAdapter) {
	go func() {
		for {
			orderWA := OrderWebAdapter{}
			ws.ReadJSON(orderWA)
			orderWAChannel <- orderWA
		}
	}()
}

func MessageSocketHandler(ws websocket.Conn, messageChannel chan packetParser.PacketUpdate) {
	go func() {
		for {
			messageWebAdapterBuf := make([]MessageWebAdapter, 100)
			timeout := time.After(time.Millisecond * 20)
		loop:
			for {
				select {
				case packet := <-messageChannel:
					adapter := newMessageWebAdapter(packet)
					messageWebAdapterBuf = append(messageWebAdapterBuf, adapter)
					if len(messageWebAdapterBuf) == 100 {
						ws.WriteJSON(messageWebAdapterBuf)
						break loop
					}
				case <-timeout:
					ws.WriteJSON(messageWebAdapterBuf)
					break loop
				}
			}
		}
	}()
}
