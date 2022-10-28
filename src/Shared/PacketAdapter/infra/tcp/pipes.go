package tcp

import (
	"net"
	"sync"
	"time"

	"github.com/HyperloopUPV-H8/Backend-H8/Shared/PacketAdapter/infra/aliases"
)

var (
	packetMaxLength int = 1024
)

type Pipes struct {
	conns map[aliases.IP]*net.TCPConn
	guard *sync.Mutex
}

func NewPipes(expectedPipes int) Pipes {
	return Pipes{
		conns: make(map[aliases.IP]*net.TCPConn, expectedPipes),
		guard: &sync.Mutex{},
	}
}

func (pipes *Pipes) AddConnection(ip aliases.IP, conn *net.TCPConn) {
	pipes.guard.Lock()
	defer pipes.guard.Unlock()

	pipes.conns[ip] = conn
}

func (pipes *Pipes) RemoveConnection(ip aliases.IP) {
	pipes.guard.Lock()
	defer pipes.guard.Unlock()

	conn, exists := pipes.conns[ip]
	if !exists {
		return
	}
	conn.Close()
	delete(pipes.conns, ip)
}

func (pipes *Pipes) Close() {
	pipes.guard.Lock()
	defer pipes.guard.Unlock()

	for _, conn := range pipes.conns {
		conn.Close()
	}
}

func (pipes *Pipes) Receive() []aliases.Payload {
	pipes.guard.Lock()
	defer pipes.guard.Unlock()

	payloads := make([]aliases.Payload, 0, len(pipes.conns))
	for _, conn := range pipes.conns {
		conn.SetDeadline(time.Now().Add(time.Second * 2))
		buf := make(aliases.Payload, packetMaxLength)
		n, _ := conn.Read(buf)
		if n > 0 {
			payloads = append(payloads, buf)
		}
	}

	return payloads
}

func (pipes *Pipes) Send(ip aliases.IP, payload aliases.Payload) {
	pipes.guard.Lock()
	defer pipes.guard.Unlock()

	conn, exists := pipes.conns[ip]
	if !exists {
		return
	}

	_, err := conn.Write(payload)

	// We use a goroutine because Close will attempt to lock but we still have the lock in Send
	if err != nil {
		go pipes.RemoveConnection(ip)
	}
}

func (pipes *Pipes) ConnectedIPs() []aliases.IP {
	pipes.guard.Lock()
	defer pipes.guard.Unlock()

	ips := make([]aliases.IP, 0, len(pipes.conns))
	for ip := range pipes.conns {
		ips = append(ips, ip)
	}

	return ips
}
