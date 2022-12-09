package transport_controller

import (
	"net"

	"github.com/HyperloopUPV-H8/Backend-H8/transport_controller/internals"
	"github.com/HyperloopUPV-H8/Backend-H8/transport_controller/models"
)

type TransportController struct {
	sniffer *internals.Sniffer
	pipes   *internals.PipeHandle
	Config  models.Config
}

func Open(laddr *net.TCPAddr, raddrs []*net.TCPAddr, device string, live bool, config models.Config) *TransportController {
	return &TransportController{
		sniffer: internals.OpenSniffer(device, live, config),
		pipes:   internals.OpenPipes(laddr, raddrs, config),
		Config:  config,
	}
}

func (controller *TransportController) Write(addr string, payload []byte) bool {
	return controller.pipes.Write(addr, payload)
}

func (controller *TransportController) Stats() (recieved int, dropped int) {
	return controller.sniffer.Stats()
}
