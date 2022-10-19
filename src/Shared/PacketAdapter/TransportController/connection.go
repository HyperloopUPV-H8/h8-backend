package TransportController

import (
	"fmt"
	"net"
	"time"

	"github.com/go-ping/ping"
)

type connection struct {
	addr    *net.TCPAddr
	tcp     *net.TCPConn
	pinger  *ping.Pinger
	isAlive bool
}

func newConnection(ip string, port int) (*connection, error) {
	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		return nil, err
	}

	pinger, err := ping.NewPinger(ip)
	if err == nil {
		// SetPrivileged is required to run on windows (should run as normal user)
		pinger.SetPrivileged(true)
		pinger.Count = 1
		pinger.Timeout = time.Second
	} else {
		return nil, err
	}

	conn, _ := net.DialTCP("tcp", nil, addr)

	return &connection{
		addr:    addr,
		pinger:  pinger,
		tcp:     conn,
		isAlive: conn != nil,
	}, nil
}

func (conn *connection) tryConnect() {
	if conn.tcp != nil {
		return
	}

	conn.pinger.Run()
	if conn.pinger.Statistics().PacketsRecv > 0 {
		conn.tcp, _ = net.DialTCP("tcp", nil, conn.addr)
		conn.isAlive = conn.tcp != nil
	}
}

func (conn *connection) checkAlive() {
	if conn.tcp == nil {
		return
	}

	_, err := conn.tcp.Write(make([]byte, 1))
	if err != nil {
		conn.disconnect()
	}
}

func (conn *connection) disconnect() {
	conn.tcp.Close()
	conn.tcp = nil
	conn.isAlive = false
}
