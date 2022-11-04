package streaming

import (
	"time"

	"github.com/HyperloopUPV-H8/Backend-H8/DataTransfer/domain"
	"github.com/gorilla/websocket"
)

func DataSocketHandler(ws websocket.Conn, packetChannel <-chan domain.Packet) {
	go func() {
		for {
			packetWebAdapterBuf := make([]PacketWebAdapter, 100)
			timeout := time.After(time.Millisecond * 20)
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
	}()
}

func OrderSocketHandler(ws websocket.Conn, orderWAChannel chan<- OrderWebAdapter) {
	go func() {
		for {
			orderWA := OrderWebAdapter{}
			ws.ReadJSON(orderWA)
			orderWAChannel <- orderWA
		}
	}()
}
