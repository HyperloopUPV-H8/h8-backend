package streaming

import (
	"time"

	"github.com/HyperloopUPV-H8/Backend-H8/DataTransfer/domain"
	"github.com/gorilla/websocket"
)

func DataSocketHandler(ws websocket.Conn, packetChannel <-chan domain.Packet) {
	for {
		packetWebAdapterBuf := make([]PacketWebAdapter, 100)
		timeout := time.After(time.Second)
	loop:
		for {
			select {
			case packet := <-packetChannel:
				adapter := newPacketWebAdapter(packet)
				packetWebAdapterBuf = append(packetWebAdapterBuf, adapter)
				if len(packetWebAdapterBuf) == 100 {
					ws.WriteJSON(packetWebAdapterBuf)
					break loop
				}
			case <-timeout:
				ws.WriteJSON(packetWebAdapterBuf)
				break loop
			}
		}
	}
}
