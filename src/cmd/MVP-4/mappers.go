package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"

	excelAdapterModels "github.com/HyperloopUPV-H8/Backend-H8/excel_adapter/models"
)

type IPtoBoard map[string]string
type IDtoIP map[uint16]string
type IDtoType map[uint16]string

func (table *IPtoBoard) AddPacket(board string, ip string, desc excelAdapterModels.Description, values []excelAdapterModels.Value) {
	if table == nil {
		*table = make(IPtoBoard)
	}

	(*table)[ip] = board
}

func (table *IDtoIP) AddPacket(board string, ip string, desc excelAdapterModels.Description, values []excelAdapterModels.Value) {
	if table == nil {
		*table = make(IDtoIP)
	}

	id, err := strconv.ParseUint(desc.ID, 10, 16)
	if err != nil {
		log.Fatalln(err)
	}
	(*table)[uint16(id)] = ip
}

func (table *IDtoType) AddPacket(board string, ip string, desc excelAdapterModels.Description, values []excelAdapterModels.Value) {
	if table == nil {
		*table = make(IDtoType)
	}

	id, err := strconv.ParseUint(desc.ID, 10, 16)
	if err != nil {
		log.Fatalln(err)
	}
	(*table)[uint16(id)] = desc.Type
}

func getFilter(raddrs []*net.TCPAddr) string {
	hosts := "("
	for _, addr := range raddrs {
		hosts += fmt.Sprintf("src host %s || ", addr.IP)
	}
	hosts = strings.TrimSuffix(hosts, " || ") + ") && ("
	for _, addr := range raddrs {
		hosts += fmt.Sprintf("dst host %s || ", addr.IP)
	}
	hosts = strings.TrimSuffix(hosts, " || ") + ")"
	return "(udp && (" + hosts + ")) || (((tcp[tcpflags] & tcp-push) > 0) && (" + hosts + "))"
}

func getJSON(data any) []byte {
	encoded, err := json.Marshal(data)
	if err != nil {
		log.Fatalln(err)
	}
	return encoded
}
