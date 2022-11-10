package tcp

import (
	"fmt"
	"log"
	"net"
	"reflect"
	"sync"
	"testing"
	"time"
)

func TestClose(t *testing.T) {
	var server Server

	t.Cleanup(func() {
		server.Close()
	})

	tests := map[string][]string{
		"single connection":    {"127.0.0.2"},
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
				if r := recover(); r != nil {
					t.Logf("test \"pipes: close: %s (should panic)\" panicked: %v", name, r)
				}
			}()

			performTCP(testAddr, fmt.Sprintf("127.0.0.1:%d", port), func(conn *net.TCPConn) {})

			performTCP(testAddr, fmt.Sprintf("127.0.0.1:%d", port), func(conn *net.TCPConn) {})

			t.Fatalf("pipes: close: %s (should panic) didn't panicked", name)
		})

		port = getTCPPort()

		testAddr = make([]string, len(test))
		for i, t := range test {
			testAddr[i] = fmt.Sprintf("%s:%d", t, port)
		}

		t.Run("pipes: close: "+name, func(t *testing.T) {
			server = Open(port, test)

			performTCP(testAddr, fmt.Sprintf("127.0.0.1:%d", port), func(conn *net.TCPConn) {})

			server.Close()
			server = Open(port, test)

			performTCP(testAddr, fmt.Sprintf("127.0.0.1:%d", port), func(conn *net.TCPConn) {})
		})
	}
}

func TestConnect(t *testing.T) {
	type testCase struct {
		addrs   []string
		connect []string
		expect  []string
	}

	tests := map[string]testCase{
		"single connection": {
			addrs:   []string{"127.0.0.2"},
			connect: []string{"127.0.0.2"},
			expect:  []string{"127.0.0.2"},
		},
		"multiple connections (all connected) (all in list)": {
			addrs:   []string{"127.0.0.2", "127.0.0.3"},
			connect: []string{"127.0.0.2", "127.0.0.3"},
			expect:  []string{"127.0.0.2", "127.0.0.3"},
		},
		"multiple connections (all in list)": {
			addrs:   []string{"127.0.0.2", "127.0.0.3", "127.0.0.4"},
			connect: []string{"127.0.0.2", "127.0.0.3"},
			expect:  []string{"127.0.0.2", "127.0.0.3"},
		},
		"multiple connections (all connected)": {
			addrs:   []string{"127.0.0.2", "127.0.0.3"},
			connect: []string{"127.0.0.2", "127.0.0.4"},
			expect:  []string{"127.0.0.2"},
		},
		"multiple connections": {
			addrs:   []string{"127.0.0.2", "127.0.0.3", "127.0.04"},
			connect: []string{"127.0.0.2", "127.0.0.5"},
			expect:  []string{"127.0.0.2"},
		},
	}

	for name, test := range tests {
		port := getTCPPort()

		testAddrs := make([]string, len(test.connect))
		for i, t := range test.connect {
			testAddrs[i] = fmt.Sprintf("%s:%d", t, port)
		}

		t.Run(name, func(t *testing.T) {
			server := Open(port, test.addrs)
			defer server.Close()

			performTCP(testAddrs, fmt.Sprintf("127.0.0.1:%d", port), func(conn *net.TCPConn) {})

			<-time.After(time.Millisecond * 50) // We need to wait because ConnectedAddresses gets updated concurrently
			got := server.ConnectedAddresses()
			if !reflect.DeepEqual(got, test.expect) {
				t.Fatalf("expected %v, got %v", test.expect, got)
			}
		})
	}
}

func BenchmarkValidAddr(b *testing.B) {
	server := Open(0, []string{"127.0.0.2", "127.0.0.3", "127.0.0.4"})

	b.Run("invalid", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			server.isValidAddr("127.0.0.1")
		}
	})

	b.Run("valid", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			server.isValidAddr("127.0.0.4")
		}
	})

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

func performTCP(laddrs []string, raddr string, action func(*net.TCPConn)) {
	conns := make([]*net.TCPConn, len(laddrs))
	defer func() {
		for _, conn := range conns {
			conn.Close()
		}
	}()

	wait := sync.WaitGroup{}

	var err error
	for i, laddr := range laddrs {
		conns[i], err = connectTCP(laddr, raddr)
		if err != nil {
			panic(err)
		}
		fmt.Println("connected", laddr)
		wait.Add(1)
		go func(conn *net.TCPConn) {
			defer wait.Done()
			action(conn)
		}(conns[i])
	}

	wait.Wait()
}

func getTCPPort() uint16 {
	conn, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Fatalln("server: get tcp port:", err)
	}
	defer conn.Close()
	return uint16(conn.Addr().(*net.TCPAddr).Port)
}
