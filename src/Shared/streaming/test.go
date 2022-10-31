package streaming

// import (
// 	"log"
// 	"net/http"
// 	"time"

// 	"github.com/gorilla/websocket"
// 	"golang.org/x/net/websocket"
// )

// // El handler tiene que ser un closure para tener acceso a cosas exteriores, en este caso DataTransfer.PacketChannel

// // Necesito ws y DataTransfer.PacketChannel

// func dataHandler(ws websocket.Conn, packetChannel <-chan Packet) {
// 	packetAdapterBuf := make([]PacketAdapter)
// 	go for packetWebAdapter := range packetChannel {
// 		webAdapter := adapters.NewPacketWebAdapter(packet)
// 		packetWebAdapter.push(webAdapter)

// 		//Meter timeout
// 		if (len(buf) > limite) {
// 			ws.send(webAdapter)
// 		}

// 	}
// }

// func getDataHandler(packetChannel) {
// 	return func(w, r) {
// 		websocket = r.upgrade()
// 		//dataHandler(websocket)
// 		orderHandle()
// 	}
// }

// func DataSocketHandler(w http.ResponseWriter, r *http.Request) {
// 	upgrader := websocket.Upgrader{}
// 	c, err := upgrader.Upgrade(w, r, nil)
// 	if err != nil {
// 		log.Println("upgrade err: ", err)
// 	}

// 	defer c.Close()

// 	buf := make([]PacketAdapter)

// 	go for packet range DataChannel {
// 		adapter := webadapters.NewPacketAdapter(packet)
// 		buf.append(adapter)

// 	select {
// 	case: time.NewTimer(time.Millisecond * 20):
// 		ws.send(buf)
// 	default:
// 		if len(buf) > 100 {
// 			ws.send(buf)
// 		}
// 	}
// 	}
// }
