package TransportController

import (
	"fmt"
	"net"
	"time"

	"github.com/go-ping/ping"
)

// Helper struct to keep track of a connection with the pod
type connection struct {
	addr    *net.TCPAddr
	tcp     *net.TCPConn
	pinger  *ping.Pinger
	isAlive bool
}

// Create a new connection and attempt to connect, returning error if the given ip and port are invalid
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

// Attempt to establish connection with the other end
func (c *connection) tryConnect() {
	if c.tcp == nil {
		c.pinger.Run()
		if c.pinger.Statistics().PacketsRecv > 0 {
			c.tcp, _ = net.DialTCP("tcp", nil, c.addr)
			c.isAlive = c.tcp != nil
		}
	}
}

// Check if the other end is still connected by sending a dumy message
func (c *connection) checkAlive() {
	if c.tcp != nil {
		_, err := c.tcp.Write(make([]byte, 1))
		if err != nil {
			c.disconnect()
		}
	}
}

// Disconnect from the other end, gracefully closing the connection
func (c *connection) disconnect() {
	c.tcp.Close()
	c.tcp = nil
	c.isAlive = false
}
