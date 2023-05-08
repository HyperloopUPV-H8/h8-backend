package models

type SortedPackets []Packet

func (packets SortedPackets) Len() int {
	return len(packets)
}

func (packets SortedPackets) Less(i, j int) bool {
	return int(packets[i].ID) < int(packets[j].ID)
}

func (packets SortedPackets) Swap(i, j int) {
	copyPacket := packets[i]
	packets[i] = packets[j]
	packets[j] = copyPacket
}
