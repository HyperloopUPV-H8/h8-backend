package dataTransfer

//"github.com/HyperloopUPV-H8/Backend-H8/src/DataTransfer/podDataCreator"
//"github.com/HyperloopUPV-H8/Backend-H8/Shared/packetParser"

type DataTransfer struct {
	pd PodData
}

func (dt DataTransfer) Invoke(getPackets func() []packetParser.PacketUpdate) {
	for {
		packets := getPackets()
		pd.UpdatePackets(packets)
	}
}
