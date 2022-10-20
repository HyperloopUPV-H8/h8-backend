package infra

import "net"

type Connection struct {
	tcp  *net.TCPConn
	addr *net.TCPAddr
}

func NewConnection(addr *net.TCPAddr) (*Connection, error) {
	// TODO
	return nil, nil
}
