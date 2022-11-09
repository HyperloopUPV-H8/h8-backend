package streaming

import (
	"time"

	"github.com/HyperloopUPV-H8/Backend-H8/data_transfer/infra/dto"
	"github.com/gorilla/websocket"
)

func DataSocketHandler(ws websocket.Conn, packetChannel <-chan dto.Packet) {
	go func() {
	routine:
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
					if err := ws.WriteJSON(packetWebAdapterBuf); err != nil {
						break routine
					}
					break loop
				}
			}
		}
	}()
}
