package PacketAdapter

import (
	"fmt"
	"net"
	"testing"
	"time"
)

// Check that NewTransportController returns a valid pointer to a TransportController all the time
func TestNewTransportController(t *testing.T) {
	type args struct {
		validAddrs []string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "empty ip list",
			args: args{
				validAddrs: []string{},
			},
		},

		{
			name: "one ip",
			args: args{
				validAddrs: []string{"127.0.0.1"},
			},
		},

		{
			name: "multiple ips",
			args: args{
				validAddrs: []string{"127.0.0.1", "127.0.0.2", "127.0.0.3"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if NewTransportController(tt.args.validAddrs) == nil {
				t.Errorf("Expected TransportController, got nil\n")
			}
		})
		serverPort++
	}
}

// Check that transport controller correctly receives all messages but blocks when there are no messages left to read
func TestTransportController_ReceiveMessage(t *testing.T) {
	controller := NewTransportController([]string{"127.0.0.1"})

	go func() {
		laddr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:5999")
		if err != nil {
			panic(err)
		}
		raddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("127.0.0.1:%d", serverPort))
		if err != nil {
			panic(err)
		}
		conn, err := net.DialTCP("tcp", laddr, raddr)
		if err != nil {
			panic(err)
		}

		for i := 0; i < 5; i++ {
			conn.Write([]byte{1, 2, 3, 4, 5})
			<-time.After(time.Millisecond)
		}
		conn.Close()
	}()

	for i := 0; i < 5; i++ {
		if got := controller.ReceiveMessage(); got == nil {
			t.Errorf("Expected []byte, got nil")
		}
	}

	done := make(chan bool)
	go func() {
		controller.ReceiveMessage()

		done <- true
	}()

	timeout := time.NewTimer(time.Second * 5)

	select {
	case <-done:
		t.Errorf("Expected blocking, got value")
	case <-timeout.C:
		break
	}

	serverPort++
}
