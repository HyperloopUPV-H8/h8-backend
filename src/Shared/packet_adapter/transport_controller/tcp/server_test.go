package tcp

import (
	"fmt"
	"log"
	"net"
	"testing"
)

func TestClose(t *testing.T) {
	var server Server
	var connections []*net.TCPConn

	cleanupConnections := func() {
		for _, conn := range connections {
			conn.Close()
		}
	}

	connect := func(addrs []string, port uint16) {
		for _, addr := range addrs {
			conn, err := connectTCP(addr, fmt.Sprintf("127.0.0.1:%d", port))
			if err != nil {
				panic("pipes: close: " + err.Error())
			}
			connections = append(connections, conn)
		}
	}

	t.Cleanup(func() {
		server.Close()
		cleanupConnections()
	})

	tests := map[string][]string{
		"simple connection":    {"127.0.0.2"},
		"multiple connections": {"127.0.0.2", "127.0.0.3"},
	}

	for name, test := range tests {
		port := getTCPPort()

		testAddr := make([]string, len(test))
		for i, t := range test {
			testAddr[i] = fmt.Sprintf("%s:%d", t, port)
		}

		t.Run("pipes: close: "+name+" (should panic)", func(t *testing.T) {
			server = Open(port, test)
			defer func() {
				server.Close()
				cleanupConnections()
				if r := recover(); r != nil {
					t.Logf("test \"pipes: close: %s (should panic)\" panicked: %v", name, r)
				}
			}()

			connect(testAddr, port)

			cleanupConnections()

			connect(testAddr, port)

			t.Fatalf("pipes: close: %s (should panic) didn't panicked", name)
		})

		port = getTCPPort()

		testAddr = make([]string, len(test))
		for i, t := range test {
			testAddr[i] = fmt.Sprintf("%s:%d", t, port)
		}

		t.Run("pipes: close: "+name, func(t *testing.T) {
			server = Open(port, test)
			defer func() {
				server.Close()
				cleanupConnections()
			}()

			connect(testAddr, port)

			server.Close()
			cleanupConnections()

			server = Open(port, test)

			connect(testAddr, port)
		})
	}
}

func connectTCP(laddr string, raddr string) (*net.TCPConn, error) {
	tcpLAddr, err := net.ResolveTCPAddr("tcp", laddr)
	if err != nil {
		return nil, err
	}

	tcpRAddr, err := net.ResolveTCPAddr("tcp", raddr)
	if err != nil {
		return nil, err
	}

	return net.DialTCP("tcp", tcpLAddr, tcpRAddr)
}

func getTCPPort() uint16 {
	conn, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Fatalln("server: get tcp port:", err)
	}
	defer conn.Close()
	return uint16(conn.Addr().(*net.TCPAddr).Port)
}
