package tcp

import (
	"fmt"
	"log"
	"net"
	"reflect"
	"sort"
	"strings"
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
		"single connection":   {"127.0.0.2"},
		"multiple connecions": {"127.0.0.2", "127.0.0.3"},
	}

	for name, test := range tests {
		port := getTCPPort()

		testAddr := make([]string, len(test))
		for i, t := range test {
			testAddr[i] = fmt.Sprintf("%s:%d", t, port)
		}

		t.Run("server: close: "+name+" (should panic)", func(t *testing.T) {
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

		t.Run("server: close: "+name, func(t *testing.T) {
			server = Open(port, test)

			performTCP(testAddr, fmt.Sprintf("127.0.0.1:%d", port), func(conn *net.TCPConn) {})

			server.Close()
			<-time.After(time.Millisecond * 50)
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

		t.Run("server: connect: "+name, func(t *testing.T) {
			server := Open(port, test.addrs)
			defer server.Close()

			performTCP(testAddrs, fmt.Sprintf("127.0.0.1:%d", port), func(conn *net.TCPConn) {})

			got := server.ConnectedAddresses()
			sort.Strings(got)
			sort.Strings(test.expect)
			if !reflect.DeepEqual(got, test.expect) {
				t.Fatalf("expected %v, got %v", test.expect, got)
			}
		})
	}
}

func TestSend(t *testing.T) {
	type testCase struct {
		addrs         []string
		payloads      map[string][]byte
		maxPayloadLen int
	}

	tests := map[string]testCase{
		"single connection": {
			addrs: []string{"127.0.0.2"},
			payloads: map[string][]byte{
				"127.0.0.2": {0xff, 0xff, 0xff},
			},
			maxPayloadLen: 3,
		},
		"multiple connections": {
			addrs: []string{"127.0.0.2", "127.0.0.3"},
			payloads: map[string][]byte{
				"127.0.0.2": {0xff, 0xff, 0xff},
				"127.0.0.3": {0xff, 0xff},
			},
			maxPayloadLen: 3,
		},
	}

	for name, test := range tests {
		port := getTCPPort()

		testAddrs := make([]string, len(test.addrs))
		for i, t := range test.addrs {
			testAddrs[i] = fmt.Sprintf("%s:%d", t, port)
		}

		t.Run("server: send: "+name, func(t *testing.T) {
			server := Open(port, test.addrs)
			serverMux := sync.Mutex{}
			defer server.Close()

			performTCP(testAddrs, fmt.Sprintf("127.0.0.1:%d", port), func(conn *net.TCPConn) {
				addr := strings.Split(conn.LocalAddr().String(), ":")[0]
				serverMux.Lock()
				server.Send(addr, test.payloads[addr])
				serverMux.Unlock()

				buffer := make([]byte, len(test.payloads[addr]))
				_, err := conn.Read(buffer)
				if err != nil {
					t.Fatal(err)
				}

				if !reflect.DeepEqual(buffer, test.payloads[addr]) {
					t.Fatalf("expected: %v, got %v", test.payloads[addr], buffer)
				}
			})
		})
	}
}

func TestDisconnect(t *testing.T) {
	port := getTCPPort()

	server := Open(port, []string{"127.0.0.2"})
	defer server.Close()

	performTCP([]string{fmt.Sprintf("127.0.0.2:%d", port)}, fmt.Sprintf("127.0.0.1:%d", port), func(conn *net.TCPConn) {
		conn.SetLinger(0)
	})

	connected := server.ConnectedAddresses()
	sort.Strings(connected)
	expect := []string{"127.0.0.2"}
	sort.Strings(expect)
	if !reflect.DeepEqual(connected, expect) {
		t.Fatalf("expected %v, got %v", expect, connected)
	}

	server.Send("127.0.0.2", []byte{0x00})

	<-time.After(time.Millisecond * 500)
	connected = server.ConnectedAddresses()
	expect = []string{}
	if !reflect.DeepEqual(connected, expect) {
		t.Fatalf("expected %v, got %v", expect, connected)
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

		wait.Add(1)
		<-time.After(time.Millisecond * 10)
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
