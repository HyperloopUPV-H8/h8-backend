package PacketAdapter

import (
	"net"
	"sync"
	"time"
)

const (
	packetMaxLength int = 1024
)

type SyncPipes struct {
	conns map[IP]net.TCPConn
	guard *sync.Mutex
}

func NewPipes(expectedPipes int) SyncPipes {
	return SyncPipes{
		conns: make(map[IP]net.TCPConn, expectedPipes),
		guard: &sync.Mutex{},
	}
}

func (pipes SyncPipes) AddConnection(ip IP, conn net.TCPConn) {
	pipes.guard.Lock()
	defer pipes.guard.Unlock()

	pipes.conns[ip] = conn
}

func (pipes SyncPipes) RemoveConnection(ip IP) {
	pipes.guard.Lock()
	defer pipes.guard.Unlock()

	conn, exists := pipes.conns[ip]
	if !exists {
		return
	}
	conn.Close()
	delete(pipes.conns, ip)
}

func (pipes SyncPipes) Close() {
	pipes.guard.Lock()
	defer pipes.guard.Unlock()

	for _, conn := range pipes.conns {
		conn.Close()
	}
}

func (pipes SyncPipes) Receive() []Payload {
	pipes.guard.Lock()
	defer pipes.guard.Unlock()

	data := make([]Payload, 0)

	for _, conn := range pipes.conns {
		conn.SetDeadline(time.Now().Add(time.Microsecond))
		buf := make(Payload, packetMaxLength)
		n, _ := conn.Read(buf)
		if n > 0 {
			data = append(data, buf)
		}
	}

	return data
}

func (pipes SyncPipes) Send(ip IP, payload Payload) error {
	pipes.guard.Lock()
	defer pipes.guard.Unlock()

	conn, exists := pipes.conns[ip]
	if !exists {
		return nil
	}

	_, err := conn.Write(payload)
	return err
}

func (pipes SyncPipes) ConnectedIPs() []IP {
	pipes.guard.Lock()
	defer pipes.guard.Unlock()

	ips := make([]IP, 0)
	for ip := range pipes.conns {
		ips = append(ips, ip)
	}

	return ips
}
