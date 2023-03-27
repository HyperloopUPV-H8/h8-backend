package models

type GlobalInfo struct {
	BoardToIP        map[string]string
	UnitToOperations map[string]string
	ProtocolToPort   map[string]string
	BoardToID        map[string]string
}
